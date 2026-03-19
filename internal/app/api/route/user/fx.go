package userroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user-route",
	common.ProvideAsRoute(NewUserRoute),
	fx.Provide(NewUserHandler),
)
