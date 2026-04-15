package mailable

type ComicNewChapterMail struct {
	*mailable
}

type ComicNewChapterMailParams struct {
	UserName           string
	ComicTitle         string
	ComicThumbnailURL  string
	ChapterDisplayName string
	ChapterNumber      string
	ChapterTitle       string
	CurrentYear        int
}

func NewComicNewChapterMail(data ComicNewChapterMailParams) MailableInterface {
	mailable := NewMailable()
	mailable.subject = "New chapter available: " + data.ComicTitle
	mailable.templateName = "comic-new-chapter"
	mailable.data = data

	return &ComicNewChapterMail{
		mailable,
	}
}
