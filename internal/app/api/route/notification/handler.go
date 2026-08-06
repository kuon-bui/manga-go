package notificationroute

import (
	"manga-go/internal/pkg/logger"
	notificationservice "manga-go/internal/pkg/services/notification"

	"go.uber.org/fx"
)

type NotificationHandler struct {
	notificationService *notificationservice.NotificationService
	logger              *logger.Logger
}

type NotificationHandlerParams struct {
	fx.In
	NotificationService *notificationservice.NotificationService
	Logger              *logger.Logger
}

func NewNotificationHandler(p NotificationHandlerParams) *NotificationHandler {
	return &NotificationHandler{
		notificationService: p.NotificationService,
		logger:              p.Logger,
	}
}
