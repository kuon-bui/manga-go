package chapterrequest

type UpdateChapterPagesRequest struct {
	Pages []PageRequest `json:"pages" binding:"required,min=1,dive"`
}
