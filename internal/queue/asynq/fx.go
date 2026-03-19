package asynq

import (
	"base-go/internal/pkg/common"
	"base-go/internal/queue/asynq/pkg/mail"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"asynq",
	fx.Provide(
		NewAsynqServer,
		NewAsynqServerMux,
	),
	common.AsWorkerManager(mail.NewMailDeliverWorker),
)
