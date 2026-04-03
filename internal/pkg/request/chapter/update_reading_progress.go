package chapterrequest

type UpdateReadingProgressRequest struct {
	ScrollPercent int `json:"scrollPercent" binding:"required,min=0,max=100" example:"50"`
}
