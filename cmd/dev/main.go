package main

import (
	"manga-go/internal/app"
	"manga-go/internal/app/api"
	"manga-go/internal/app/api/server"
	"manga-go/internal/pkg/tracer"

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
