package handler

import (
	"context"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/model"
	"github.com/thetkpark/uob-thai-credit-card-notification-common/publisher"
	"log"
	"log/slog"
	"regexp"
)

type Handler struct {
	publisher      publisher.Publisher
	whiteListCards []string
	usageRegex     *regexp.Regexp
}

func NewHandlerImpl(publisher publisher.Publisher, whiteListCards []string) *Handler {
	usageNotiRegex, err := regexp.Compile("^มีการใช้บัตร UOB-(\\d{4}) @(.+) (\\d+\\.\\d{2} THB) วันที่ (\\d{2}/\\d{2})")
	if err != nil {
		log.Fatal(err)
	}

	return &Handler{
		publisher:      publisher,
		whiteListCards: whiteListCards,
		usageRegex:     usageNotiRegex,
	}
}

func (h Handler) HandleUsageNotificationText(ctx context.Context, text string) {
	slog.DebugContext(ctx, "Received message", slog.String("message", text))

	if !h.usageRegex.MatchString(text) {
		slog.InfoContext(ctx, "Not a usage notification", slog.String("message", text))
		return
	}
	match := h.usageRegex.FindStringSubmatch(text)
	if len(match) != 5 {
		slog.ErrorContext(ctx, "Invalid usage notification", slog.String("message", text))
		return
	}

	if !h.isCardInWhiteList(match[1]) {
		slog.InfoContext(ctx, "Card not in white list", slog.String("card", match[1]))
		return
	}

	noti := model.UsageNotification{
		Shop:   match[2],
		Amount: match[3],
		Date:   match[4],
	}

	if err := h.publisher.PublishMessage(ctx, noti); err != nil {
		slog.ErrorContext(ctx, "Error publishing message", slog.String("error", err.Error()))
		return
	}

	slog.InfoContext(ctx, "Usage notification published", slog.Any("notification", noti))
}

func (h Handler) isCardInWhiteList(cardNumber string) bool {
	for _, c := range h.whiteListCards {
		if c == cardNumber {
			return true
		}
	}
	return false
}
