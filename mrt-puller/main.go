package main

import (
	"context"
	"fmt"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/kv"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/logger"
	commonmodel "github.com/thetkpark/uob-thai-credit-card-notification-common/model"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/publisher"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/trace"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/mrt-puller/config"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/mrt-puller/model"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/mrt-puller/mrt"
	"log"
	"log/slog"
)

const LatestJourneyKey = "latest_journey"

func main() {
	conf := config.Init()
	logger.Init(conf.Log, true)

	m := mrt.NewApiImpl(conf.Mrt)
	kv := kv.NewRedisKV(conf.KV)
	pub := publisher.NewPubSubPublisher(conf.Publisher.ProjectID, conf.Publisher.TopicID)

	correlationId := trace.GenerateCorrelationId()
	ctx := trace.AddCorrelationIdToLogContext(context.Background(), correlationId)
	ctx = context.WithValue(ctx, trace.CorrelationIdKey, correlationId)

	latestProcessJourneyID, err := kv.Get(ctx, LatestJourneyKey)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting latest journey id", err)
	}

	accessToken, err := m.GetAccessToken(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting access token", err)
		log.Fatalln("Error getting access token")
	}

	req := model.GetJourneyRequest{
		CardID:      conf.CardId,
		PageNo:      1,
		PageSize:    conf.Mrt.FetchLimit,
		AccessToken: accessToken,
		Lang:        model.LangTH,
	}
	journeyResp, err := m.GetJourney(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting journey", err)
		log.Fatalln("Error getting journey")
	}

	var journeys []model.JourneyData

out:
	for _, journeyDate := range journeyResp.Data.List {
		for _, j := range journeyDate.Journeys {
			if j.JourneyID == latestProcessJourneyID {
				break out
			}
			if exist, err := kv.Exist(ctx, j.JourneyID); err != nil {
				slog.ErrorContext(ctx, "Error checking existance", slog.String("error", err.Error()))
			} else if exist {
				continue
			}
			journeys = append(journeys, j)
		}
	}

	if len(journeys) == 0 {
		slog.InfoContext(ctx, "No new journey")
		return
	}

	for _, j := range journeys {
		msg := commonmodel.UsageNotification{
			Shop:   fmt.Sprintf("MRT %s -> %s", j.From.StationName, j.To.StationName),
			Amount: fmt.Sprintf("%d THB", j.TotalAmount),
			Date:   j.Date,
		}
		slog.Debug("Message", slog.Any("message", msg))

		if err := pub.PublishMessage(ctx, msg); err != nil {
			slog.ErrorContext(ctx, "Error publishing message",
				slog.String("error", err.Error()),
				slog.String("message", fmt.Sprintf("%+v", msg)),
				slog.String("journey_id", j.JourneyID),
			)
		}
	}

	// set latest journey id
	if err := kv.Add(ctx, LatestJourneyKey, journeys[0].JourneyID, 0); err != nil {
		slog.ErrorContext(ctx, "Error setting latest journey id", slog.String("error", err.Error()))
		log.Fatalln("Error setting latest journey id")
	}
}
