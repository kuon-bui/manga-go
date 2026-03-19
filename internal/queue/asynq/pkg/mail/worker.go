package mail

import (
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/mail"
	queueconstant "manga-go/internal/queue/queue_constant"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type MailDeliverParams struct {
	fx.In

	Mux        *asynq.ServeMux
	Logger     *logger.Logger
	MailDialer *mail.MailDialer
}

type MailDeliverWorker struct {
	mux        *asynq.ServeMux
	logger     *logger.Logger
	mailDialer *mail.MailDialer
}

func NewMailDeliverWorker(p MailDeliverParams) *MailDeliverWorker {
	return &MailDeliverWorker{
		mux:        p.Mux,
		logger:     p.Logger,
		mailDialer: p.MailDialer,
	}
}

func (w *MailDeliverWorker) RegisterWorkers() {
	w.mux.HandleFunc(queueconstant.MAIL_DELIVER_TASK, w.mailDeliverHandler)
	w.mux.HandleFunc(queueconstant.MULTI_MAIL_DELIVER_TASK, w.multiMailDeliverHandler)
}
