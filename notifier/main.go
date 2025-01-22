package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/logger"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/notifier/config"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/notifier/handler"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/notifier/model"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/notifier/notification"
	"log/slog"
)

type MessagePublishedData struct {
	Message PubSubMessage
}
type PubSubMessage struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

func init() {
	c := config.Init()
	logger.Init(c.Log, true)
	noti := notification.NewTelegramNotification(c.TelegramBot.Token, c.TelegramBot.TargetChatIds)
	h := handler.NewHandlerImpl(noti)

	functions.CloudEvent("Notifier", func(ctx context.Context, e v2.Event) error {
		var msg MessagePublishedData
		if err := e.DataAs(&msg); err != nil {
			slog.ErrorContext(ctx, "Error parsing data", slog.String("error", err.Error()))
			return nil
		}

		if msg.Message.Attributes == nil || msg.Message.Attributes["x-correlation-id"] == "" {
			ctx = logger.AppendCtxValue(ctx, slog.String("x-correlation-id", uuid.NewString()))
		} else {
			ctx = logger.AppendCtxValue(ctx, slog.String("x-correlation-id", msg.Message.Attributes["x-correlation-id"]))
		}

		slog.InfoContext(ctx, fmt.Sprintf("Received message: %s", e.ID()))

		var usage model.UsageNotification
		if err := json.Unmarshal(msg.Message.Data, &usage); err != nil {
			slog.ErrorContext(ctx, "Error unmarshalling data", slog.String("error", err.Error()), slog.String("data", string(msg.Message.Data)))
			return nil
		}
		if err := h.HandleUsageNotification(ctx, usage); err != nil {
			slog.ErrorContext(ctx, "Error handling usage notification", slog.String("error", err.Error()), slog.Group("usage", usage))
			return err
		}

		return nil
	})

	slog.Info("Service started")
}
