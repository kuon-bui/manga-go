package notificationroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"notification-route",
	common.ProvideAsRoute(NewNotificationRoute),
	fx.Provide(NewNotificationHandler),
)
