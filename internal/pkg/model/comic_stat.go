package model

import (
	"time"

	"github.com/google/uuid"
)

type ComicStat struct {
	ComicID      uuid.UUID `gorm:"column:comic_id;primaryKey"`
	FollowCount  int       `gorm:"column:follow_count"`
	RatingCount  int       `gorm:"column:rating_count"`
	ChapterCount int       `gorm:"column:chapter_count"`
	AvgRating    *float64  `gorm:"column:avg_rating"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (ComicStat) TableName() string {
	return "comic_stats"
}
