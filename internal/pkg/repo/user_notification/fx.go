package usernotificationrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"user-notification-repo",
	fx.Provide(NewUserNotificationRepo),
)
