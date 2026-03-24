package chapterrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type ChapterRepo struct {
	*base.BaseRepository[model.Chapter]
}

func NewChapterRepo(db *gorm.DB) *ChapterRepo {
	return &ChapterRepo{
		BaseRepository: &base.BaseRepository[model.Chapter]{
			DB: db,
		},
	}
}
