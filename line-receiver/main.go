package main

import (
	"context"
	"errors"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/logger"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/publisher"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/line-receiver/config"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/line-receiver/handler"
	"log"
	"net/http"
)

type UsageNotification struct {
	Shop   string `json:"shop"`
	Amount string `json:"amount"`
	Date   string `json:"date"`
}

func main() {
	c := config.Init()
	logger.Init(c.Log, true)

	pub := publisher.NewPubSubPublisher(c.PubSub.ProjectID, c.PubSub.TopicID)
	h := handler.NewHandlerImpl(pub, c.WhiteListCardNumbers)

	http.HandleFunc("/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		cb, err := webhook.ParseRequest(c.LineChannelSecret, r)
		if err != nil {
			log.Printf("Cannot parse request: %+v\n", err)
			if errors.Is(err, webhook.ErrInvalidSignature) {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range cb.Events {
			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					h.HandleUsageNotificationText(context.Background(), message.Text)
				}
			}
		}
	},
	)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
