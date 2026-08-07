package authorizationroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"authorization-route",
	common.ProvideAsRoute(NewRoute),
	fx.Provide(NewHandler),
)
