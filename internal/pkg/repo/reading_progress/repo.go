package readingprogressrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type ReadingProgressRepo struct {
	*base.BaseRepository[model.ReadingProgress]
}

func NewReadingProgressRepo(db *gorm.DB) *ReadingProgressRepo {
	return &ReadingProgressRepo{
		&base.BaseRepository[model.ReadingProgress]{
			DB: db,
		},
	}
}
