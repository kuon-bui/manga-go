package mailable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"manga-go/internal/pkg/utils"
	queueconstant "manga-go/internal/queue/queue_constant"

	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"github.com/wneessen/go-mail"
)

type MailableInterface interface {
	AddTo(address ...string) MailableInterface
	AddToFormat(mailAddress ...MailAddress) MailableInterface

	AddCc(address ...string) MailableInterface
	AddCcFormat(mailAddress ...MailAddress) MailableInterface

	AddBcc(address ...string) MailableInterface
	AddBccFormat(mailAddress ...MailAddress) MailableInterface

	SetSubject(subject string) MailableInterface
	AddAttachmentFile(file string) MailableInterface
	SetContent(content string) MailableInterface
	AddAttachmentReader(name string, data io.Reader) MailableInterface

	Build() (*mail.Msg, error)
	Dispatch(asynqClient *asynq.Client) error

	CreateMailablePayload() *MailablePayload
}

type MailAddress struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
type attachmentReader struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type mailable struct {
	to                []MailAddress
	cc                []MailAddress
	bcc               []MailAddress
	subject           string
	attachmentReaders []attachmentReader
	attachFiles       []string
	data              interface{}
	templateName      string
	content           string
}

type MailablePayload struct {
	To                []MailAddress      `json:"to"`
	Cc                []MailAddress      `json:"cc"`
	Bcc               []MailAddress      `json:"bcc"`
	Subject           string             `json:"subject"`
	AttachmentReaders []attachmentReader `json:"attachmentReaders"`
	AttachFiles       []string           `json:"attachFiles"`
	Data              interface{}        `json:"data"`
	TemplateName      string             `json:"templateName"`
	Content           string             `json:"content"`
}

var _ MailableInterface = &mailable{}

func NewMailable() *mailable {
	return &mailable{
		to:                make([]MailAddress, 0),
		cc:                make([]MailAddress, 0),
		bcc:               make([]MailAddress, 0),
		subject:           "",
		attachmentReaders: make([]attachmentReader, 0),
		attachFiles:       make([]string, 0),
		data:              nil,
		templateName:      "",
		content:           "",
	}
}

func NewMailFromPayload(data *MailablePayload) MailableInterface {
	mailable := NewMailable()
	mailable.to = data.To
	mailable.cc = data.Cc
	mailable.bcc = data.Bcc
	mailable.subject = data.Subject
	mailable.attachmentReaders = data.AttachmentReaders
	mailable.attachFiles = data.AttachFiles
	mailable.data = data.Data
	mailable.templateName = data.TemplateName

	return mailable
}

func (m *mailable) AddTo(address ...string) MailableInterface {
	m.to = append(m.to, lo.Reduce(address, func(res []MailAddress, ad string, _ int) []MailAddress {
		return append(res, MailAddress{Name: "", Address: ad})
	}, []MailAddress{})...)

	return m
}

func (m *mailable) AddToFormat(mailAddress ...MailAddress) MailableInterface {
	m.to = append(m.to, mailAddress...)
	return m
}

func (m *mailable) AddBccFormat(mailAddress ...MailAddress) MailableInterface {
	m.bcc = append(m.bcc, mailAddress...)
	return m
}

func (m *mailable) AddBcc(address ...string) MailableInterface {
	m.bcc = append(m.bcc, lo.Reduce(address, func(res []MailAddress, ad string, _ int) []MailAddress {
		return append(res, MailAddress{Name: "", Address: ad})
	}, []MailAddress{})...)

	return m
}

func (m *mailable) AddCcFormat(mailAddress ...MailAddress) MailableInterface {
	m.cc = append(m.cc, mailAddress...)
	return m
}

func (m *mailable) AddCc(address ...string) MailableInterface {
	m.cc = append(m.cc, lo.Reduce(address, func(res []MailAddress, ad string, _ int) []MailAddress {
		return append(res, MailAddress{Name: "", Address: ad})
	}, []MailAddress{})...)

	return m
}

func (m *mailable) SetSubject(subject string) MailableInterface {
	m.subject = subject
	return m
}

func (m *mailable) AddAttachmentFile(file string) MailableInterface {
	m.attachFiles = append(m.attachFiles, file)
	return m
}

func (m *mailable) AddAttachmentReader(name string, data io.Reader) MailableInterface {
	var bytes bytes.Buffer
	io.Copy(&bytes, data)
	m.attachmentReaders = append(m.attachmentReaders, attachmentReader{Name: name, Data: bytes.Bytes()})
	return m
}

func (m *mailable) SetContent(content string) MailableInterface {
	m.content = content
	return m
}

func (m *mailable) Dispatch(asynqClient *asynq.Client) error {
	payload, err := json.Marshal(m.CreateMailablePayload())
	if err != nil {
		return err
	}

	t := asynq.NewTask(queueconstant.MAIL_DELIVER_TASK, payload, asynq.MaxRetry(5))
	_, err = asynqClient.Enqueue(t, asynq.Queue(queueconstant.MAIL_DELIVER_QUEUE))

	return err
}

func (m *mailable) Build() (*mail.Msg, error) {
	msg := mail.NewMsg()
	msg.SetDate()
	msg.SetBulk()
	msg.Subject(m.subject)
	for _, to := range m.to {
		if len(to.Name) > 0 {
			msg.AddToFormat(to.Name, to.Address)
		} else {
			msg.AddTo(to.Address)
		}
	}

	for _, bcc := range m.bcc {
		if len(bcc.Name) > 0 {
			msg.AddBccFormat(bcc.Name, bcc.Address)
		} else {
			msg.AddBcc(bcc.Address)
		}
	}

	for _, cc := range m.cc {
		if len(cc.Name) > 0 {
			msg.AddCcFormat(cc.Name, cc.Address)
		} else {
			msg.AddCc(cc.Address)
		}
	}

	for _, attachFile := range m.attachFiles {
		msg.AttachFile(attachFile)
	}

	for _, attachReader := range m.attachmentReaders {
		msg.AttachReader(attachReader.Name, bytes.NewReader(attachReader.Data))
	}

	if m.templateName != "" {
		template, ok := utils.GetMailTemplate(m.templateName)
		if !ok {
			return nil, fmt.Errorf("mail template %s not found", m.templateName)
		}

		msg.AddAlternativeHTMLTemplate(template, m.data)
	} else {
		msg.SetBodyString(mail.TypeTextPlain, m.content)
	}

	return msg, nil
}

func (m *mailable) CreateMailablePayload() *MailablePayload {
	return &MailablePayload{
		To:                m.to,
		Cc:                m.cc,
		Bcc:               m.bcc,
		Subject:           m.subject,
		AttachmentReaders: m.attachmentReaders,
		AttachFiles:       m.attachFiles,
		Data:              m.data,
		TemplateName:      m.templateName,
		Content:           m.content,
	}
}
