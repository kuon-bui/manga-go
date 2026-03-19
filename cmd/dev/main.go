package main

import (
	"base-go/internal/app"
	"base-go/internal/app/api"
	"base-go/internal/app/api/server"
	"base-go/internal/pkg/tracer"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.Module,
		api.Module,
		tracer.Module,
		fx.Invoke(server.RunServer),
	).Run()
}
