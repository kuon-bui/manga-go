package notificationroute

import (
	notificationservice "manga-go/internal/pkg/services/notification"

	"go.uber.org/fx"
)

type NotificationHandler struct {
	notificationService *notificationservice.NotificationService
}

type NotificationHandlerParams struct {
	fx.In
	NotificationService *notificationservice.NotificationService
}

func NewNotificationHandler(p NotificationHandlerParams) *NotificationHandler {
	return &NotificationHandler{notificationService: p.NotificationService}
}
