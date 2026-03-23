package model

import (
	"time"

	"github.com/google/uuid"
)

type ComicAuthor struct {
	ComicID   uuid.UUID `json:"comicId" gorm:"column:comic_id;primaryKey"`
	AuthorID  uuid.UUID `json:"authorId" gorm:"column:author_id;primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (ComicAuthor) TableName() string {
	return "comic_authors"
}
