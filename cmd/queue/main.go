package main

import (
	asynqserver "manga-go/internal/app/asynq"
	"manga-go/internal/pkg/tracer"
	"manga-go/internal/queue"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		tracer.Module,
		queue.Module,
		fx.Invoke(asynqserver.RunServer),
	).Run()
}
