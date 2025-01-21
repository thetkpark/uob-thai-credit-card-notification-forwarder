package handler

import (
	"context"
	"fmt"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/notifier/model"
	"github.com/thetkpark/uob-thai-credit-card-notification-forwarder/notifier/notification"
)

type Handler interface {
	HandleUsageNotification(req model.UsageNotification) error
}

type Impl struct {
	noti notification.Notification
}

func NewHandlerImpl(noti notification.Notification) *Impl {
	return &Impl{noti: noti}
}

func (h *Impl) HandleUsageNotification(ctx context.Context, req model.UsageNotification) error {
	msg := fmt.Sprintf("Date: %s\nShop: %s\nAmount: %s", req.Shop, req.Amount, req.Date)
	if err := h.noti.BroadcastNotification(ctx, msg); err != nil {
		return err
	}
	return nil
}
