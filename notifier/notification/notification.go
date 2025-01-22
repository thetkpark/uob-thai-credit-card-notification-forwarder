package notification

import "context"

type Notification interface {
	BroadcastNotification(ctx context.Context, msg string) error
}
