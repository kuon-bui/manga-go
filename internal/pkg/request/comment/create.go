package commentrequest

import "github.com/google/uuid"

type CreateCommentRequest struct {
	ChapterID uuid.UUID  `json:"chapterId" binding:"required"`
	ParentId  *uuid.UUID `json:"parentId,omitempty"`
	Content   string     `json:"content" binding:"required"`
	PageIndex *int       `json:"pageIndex,omitempty"`
}
