package chapterrequest

type CreateChapterRequest struct {
	Number string        `json:"number" binding:"required,max=10"`
	Title  string        `json:"title" binding:"required"`
	Slug   string        `json:"slug"` // Optional - auto-generated from title if not provided
	Pages  []PageRequest `json:"pages" binding:"required,min=1,dive"`
}
