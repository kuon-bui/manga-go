package chapterrequest

type UpdateChapterRequest struct {
	Number string `json:"number" binding:"required,max=10"`
	Title  string `json:"title" binding:"required"`
	Slug   string `json:"slug" binding:"required"`
}
