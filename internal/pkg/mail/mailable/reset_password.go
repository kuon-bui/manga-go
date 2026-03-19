package mailable

type ResetPasswordMail struct {
	*mailable
}

type ResetPasswordMailParams struct {
	UserName         string
	ResetPasswordURL string
	ExpiryMinutes    int
	CurrentYear      int
	FacebookURL      string
	LinkedInURL      string
	PinterestURL     string
	YouTubeURL       string
}

func NewResetPasswordMail(data ResetPasswordMailParams) MailableInterface {
	mailable := NewMailable()
	mailable.subject = "Reset your password"
	mailable.templateName = "reset-password"
	mailable.data = data

	return &ResetPasswordMail{
		mailable,
	}
}
