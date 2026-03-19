package api

import (
	"base-go/internal/app/api/route"
	"base-go/internal/app/api/server"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"api",
	route.Module,
	fx.Provide(
		server.NewGinEngine,
		server.NewHttpServer,
	),
)
