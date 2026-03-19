package queue

import (
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/mail"
	"manga-go/internal/pkg/tracer"
	"manga-go/internal/queue/asynq"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"queue",
	fx.Provide(
		config.LoadConfig,
		tracer.InitTracer,
		logger.NewLogger,
		mail.NewMailDialer,
	),
	asynq.Module,
)
