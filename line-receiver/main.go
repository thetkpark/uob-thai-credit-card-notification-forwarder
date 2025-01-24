package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/kv"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/logger"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/publisher"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/line-receiver/config"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/line-receiver/handler"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/line-receiver/middleware"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type UsageNotification struct {
	Shop   string `json:"shop"`
	Amount string `json:"amount"`
	Date   string `json:"date"`
}

func main() {
	conf := config.Init()
	logger.Init(conf.Log, true)
	if conf.Log.ENV != "local" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	pub := publisher.NewPubSubPublisher(conf.PubSub.ProjectID, conf.PubSub.TopicID)
	k := kv.NewRedisKV(conf.Redis)
	h := handler.NewHandlerImpl(pub, k, conf.WhiteListCardNumbers)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.AttachCorrelationID())
	r.Use(middleware.HttpLogger())

	r.GET("/healthy", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	r.POST("/webhook", func(c *gin.Context) {
		cb, err := webhook.ParseRequest(conf.LineChannelSecret, c.Request)
		if err != nil {
			slog.ErrorContext(c, "Cannot parse request", slog.String("error", err.Error()))
			if errors.Is(err, webhook.ErrInvalidSignature) {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		for _, event := range cb.Events {
			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					h.HandleUsageNotificationText(c.Request.Context(), message.Text)
				}
			}
		}
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	slog.Info("Server started")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		slog.Info("timeout of 5 seconds.")
	}
	slog.Info("Server exiting")

}
