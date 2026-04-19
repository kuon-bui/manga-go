package commentrequest

import "github.com/google/uuid"

type CreateCommentRequest struct {
	// Comic-level comment: provide comicId only
	// Chapter/Page-level comment: provide chapterId (comicId inferred from chapter)
	ComicID   *uuid.UUID `json:"comicId,omitempty"`
	ChapterID *uuid.UUID `json:"chapterId,omitempty"`
	ParentId  *uuid.UUID `json:"parentId,omitempty"`
	Content   string     `json:"content" binding:"required"`
	PageIndex *int       `json:"pageIndex,omitempty"`
}
