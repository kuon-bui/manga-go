package mail

import (
	"base-go/internal/pkg/config"
	"base-go/internal/pkg/logger"
	"base-go/internal/pkg/mail/mailable"
	"base-go/internal/pkg/utils"
	queueconstant "base-go/internal/queue/queue_constant"
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"github.com/wneessen/go-mail"
	"go.uber.org/fx"
)

type MailDialer struct {
	client *mail.Client
	config *config.Config
	logger *logger.Logger
}

type mailDialerParams struct {
	fx.In

	Config *config.Config
	Logger *logger.Logger
}

func NewMailDialer(p mailDialerParams) *MailDialer {
	utils.LoadMailTemplate(p.Logger)
	var (
		tlsPolicy mail.TLSPolicy    = mail.NoTLS
		mailAuth  mail.SMTPAuthType = mail.SMTPAuthNoAuth
		username                    = p.Config.SMTP.Username
		password                    = p.Config.SMTP.Password
	)

	if len(username)+len(password) > 0 {
		tlsPolicy = mail.DefaultTLSPolicy
		mailAuth = mail.SMTPAuthAutoDiscover
	}

	client, err := mail.NewClient(
		p.Config.SMTP.Host,
		mail.WithPort(p.Config.SMTP.Port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTLSPortPolicy(tlsPolicy),
	)

	client.SetDebugLog(p.Config.Service.DebugMode)
	client.SetSMTPAuth(mailAuth)

	if err != nil {
		p.Logger.Errorf("error creating client: %v", err)
		panic(err)
	}

	return &MailDialer{
		client: client,
		config: p.Config,
		logger: p.Logger,
	}
}

func (m *MailDialer) Send(msgs ...*mail.Msg) error {
	msgs = lo.Map(msgs, func(ms *mail.Msg, _ int) *mail.Msg {
		ms.FromFormat(m.config.SMTP.FromName, m.config.SMTP.FromMail)
		return ms
	})

	err := m.client.DialAndSend(msgs...)
	if err != nil {
		m.logger.Errorf("error sending mail: %v", err)
	}

	return err
}

func (m *MailDialer) Dispatch(asynqClient *asynq.Client, mails ...mailable.MailableInterface) error {
	payload := lo.Reduce(
		mails,
		func(res []*mailable.MailablePayload, el mailable.MailableInterface, _ int) []*mailable.MailablePayload {
			return append(res, el.CreateMailablePayload())
		},
		[]*mailable.MailablePayload{},
	)

	payloadSend, err := json.Marshal(payload)

	if err != nil {
		m.logger.Errorf("Error marshaling payload: %v", err)
		return err
	}

	t := asynq.NewTask(queueconstant.MULTI_MAIL_DELIVER_TASK, payloadSend, asynq.MaxRetry(5))
	_, err = asynqClient.Enqueue(t, asynq.Queue(queueconstant.MAIL_DELIVER_QUEUE))

	return err
}

func (m *MailDialer) SendWithContext(ctx context.Context, msgs ...*mail.Msg) error {
	msgs = lo.Map(msgs, func(ms *mail.Msg, _ int) *mail.Msg {
		ms.FromFormat(m.config.SMTP.FromName, m.config.SMTP.FromMail)
		return ms
	})

	err := m.client.DialAndSendWithContext(ctx, msgs...)
	if err != nil {
		m.logger.Errorf("error sending mail: %v", err)
	}

	return err
}
