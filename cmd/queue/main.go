package main

import (
	asynqserver "base-go/internal/app/asynq"
	"base-go/internal/pkg/tracer"
	"base-go/internal/queue"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		tracer.Module,
		queue.Module,
		fx.Invoke(asynqserver.RunServer),
	).Run()
}
