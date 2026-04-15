package notificationservice

import "go.uber.org/fx"

var Module = fx.Module(
	"notification-service",
	fx.Provide(NewNotificationService),
)
