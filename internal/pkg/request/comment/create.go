package commentrequest

import "github.com/google/uuid"

type CreateCommentRequest struct {
	ChapterID uuid.UUID `json:"chapterId" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	PageIndex *int      `json:"pageIndex,omitempty"`
}
