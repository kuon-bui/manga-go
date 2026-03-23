package model

import (
	"time"

	"github.com/google/uuid"
)

type ComicTag struct {
	ComicID   uuid.UUID `json:"comicId" gorm:"column:comic_id;primaryKey"`
	TagID     uuid.UUID `json:"tagId" gorm:"column:tag_id;primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (ComicTag) TableName() string {
	return "comic_tags"
}
