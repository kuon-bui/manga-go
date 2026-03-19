package userroute

import (
	"base-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user-route",
	common.ProvideAsRoute(NewUserRoute),
	fx.Provide(NewUserHandler),
)
