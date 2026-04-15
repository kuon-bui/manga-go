package notification

import (
	"manga-go/internal/pkg/logger"
	notificationservice "manga-go/internal/pkg/services/notification"
	queueconstant "manga-go/internal/queue/queue_constant"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type NotificationWorkerParams struct {
	fx.In
	Mux                 *asynq.ServeMux
	Logger              *logger.Logger
	NotificationService *notificationservice.NotificationService
}

type NotificationWorker struct {
	mux                 *asynq.ServeMux
	logger              *logger.Logger
	notificationService *notificationservice.NotificationService
}

func NewNotificationWorker(p NotificationWorkerParams) *NotificationWorker {
	return &NotificationWorker{
		mux:                 p.Mux,
		logger:              p.Logger,
		notificationService: p.NotificationService,
	}
}

func (w *NotificationWorker) RegisterWorkers() {
	w.mux.HandleFunc(queueconstant.NOTIFICATION_FANOUT_TASK, w.notificationFanoutHandler)
}
