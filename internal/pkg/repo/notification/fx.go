package notificationrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"notification-repo",
	fx.Provide(NewNotificationRepo),
)
