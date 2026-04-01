package chapterrequest

type UpdateChapterPagesRequest struct {
	Pages []string `json:"pages" binding:"required,min=1,dive,required"`
}
