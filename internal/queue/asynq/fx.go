package asynq

import (
	"manga-go/internal/pkg/common"
	imageworker "manga-go/internal/queue/asynq/pkg/image"
	"manga-go/internal/queue/asynq/pkg/mail"
	notificationworker "manga-go/internal/queue/asynq/pkg/notification"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"asynq",
	fx.Provide(
		NewAsynqServer,
		NewAsynqServerMux,
	),
	common.AsWorkerManager(mail.NewMailDeliverWorker),
	common.AsWorkerManager(notificationworker.NewNotificationWorker),
	common.AsWorkerManager(imageworker.NewImageProcessWorker),
)
