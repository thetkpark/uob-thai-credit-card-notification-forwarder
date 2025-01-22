package notification

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"log/slog"
)

type TelegramNotification struct {
	b             *bot.Bot
	targetChatIds []string
}

func NewTelegramNotification(token string, targetChatIds []string) *TelegramNotification {
	b, err := bot.New(token, bot.WithDefaultHandler(botDefaultHandler))
	if err != nil {
		slog.Error("Error creating telegram bot", slog.String("error", err.Error()))
		log.Fatalln("Error creating telegram bot")
	}
	b.Start(context.Background())
	return &TelegramNotification{b: b, targetChatIds: targetChatIds}
}

func botDefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.Text != "/echo" {
		return
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Your chat id is: %d", update.Message.Chat.ID),
	})
	if err != nil {
		slog.Error("Error sending message", slog.String("error", err.Error()))
	}
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
