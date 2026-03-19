package queue

import (
	"base-go/internal/pkg/config"
	"base-go/internal/pkg/logger"
	"base-go/internal/pkg/mail"
	"base-go/internal/pkg/tracer"
	"base-go/internal/queue/asynq"

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
