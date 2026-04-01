package readinghistoryrequest

import "github.com/google/uuid"

type CreateReadingHistoryRequest struct {
	ChapterID uuid.UUID `json:"chapterId" binding:"required"`
	ComicID   uuid.UUID `json:"comicId" binding:"required"`
}
