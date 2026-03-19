package mailable

type HelloMail struct {
	*mailable
}

type HelloMailParams struct {
	Title string
	Name  string
}

func NewHelloMail(data HelloMailParams) MailableInterface {
	mailable := NewMailable()
	mailable.subject = "Hello mail"
	mailable.templateName = "hello.html"
	mailable.data = data

	return &HelloMail{
		mailable,
	}
}
