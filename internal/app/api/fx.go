package api

import (
	"manga-go/internal/app/api/route"
	"manga-go/internal/app/api/server"

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
