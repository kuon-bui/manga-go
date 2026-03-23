package model

import (
	"time"

	"github.com/google/uuid"
)

type ComicGenre struct {
	ComicID   uuid.UUID `json:"comicId" gorm:"column:comic_id;primaryKey"`
	GenreID   uuid.UUID `json:"genreId" gorm:"column:genre_id;primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (ComicGenre) TableName() string {
	return "comic_genres"
}
