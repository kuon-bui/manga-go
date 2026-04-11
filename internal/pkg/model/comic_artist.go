package model

import (
	"time"

	"github.com/google/uuid"
)

type ComicArtist struct {
	ComicID   uuid.UUID `json:"comicId" gorm:"column:comic_id;primaryKey"`
	ArtistID  uuid.UUID `json:"artistId" gorm:"column:artist_id;primaryKey"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (ComicArtist) TableName() string {
	return "comic_artists"
}
