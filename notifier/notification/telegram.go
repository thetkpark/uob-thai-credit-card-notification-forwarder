package notification

import (
	"context"
	"github.com/go-telegram/bot"
	"log"
	"log/slog"
)

type TelegramNotification struct {
	b             *bot.Bot
	targetChatIds []string
}

func NewTelegramNotification(token string, targetChatIds []string) *TelegramNotification {
	b, err := bot.New(token)
	if err != nil {
		slog.Error("Error creating telegram bot", slog.String("error", err.Error()))
		log.Fatalln("Error creating telegram bot")
	}
	return &TelegramNotification{b: b, targetChatIds: targetChatIds}
}

func (t TelegramNotification) BroadcastNotification(ctx context.Context, msg string) error {
	for _, to := range t.targetChatIds {
		p := &bot.SendMessageParams{
			ChatID: to,
			Text:   msg,
		}
		if _, err := t.b.SendMessage(ctx, p); err != nil {
			return err
		}
	}
	return nil
}
