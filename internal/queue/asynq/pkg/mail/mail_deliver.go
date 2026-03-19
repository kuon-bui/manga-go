package mail

import (
	"base-go/internal/pkg/mail/mailable"
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/wneessen/go-mail"
)

type MailDeliverPayload struct {
}

func (w *MailDeliverWorker) mailDeliverHandler(ctx context.Context, t *asynq.Task) error {
	var payload mailable.MailablePayload
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		w.logger.Error("Failed to unmarshal payload: ", err)
		return err
	}

	mail := mailable.NewMailFromPayload(&payload)
	msg, err := mail.Build()
	if err != nil {
		return err
	}

	return w.mailDialer.SendWithContext(ctx, msg)
}

func (w *MailDeliverWorker) multiMailDeliverHandler(ctx context.Context, t *asynq.Task) error {
	var payload []mailable.MailablePayload
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		w.logger.Error("Failed to unmarshal payload: ", err)
		return err
	}
	mails := []*mail.Msg{}
	for _, m := range payload {
		mail := mailable.NewMailFromPayload(&m)
		msg, err := mail.Build()
		if err != nil {
			return err
		}

		mails = append(mails, msg)
	}

	return w.mailDialer.SendWithContext(ctx, mails...)
}
