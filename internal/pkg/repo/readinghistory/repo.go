package readinghistoryrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type ReadingHistoryRepo struct {
	*base.BaseRepository[model.ReadingHistory]
}

func NewReadingHistoryRepo(db *gorm.DB) *ReadingHistoryRepo {
	return &ReadingHistoryRepo{
		BaseRepository: &base.BaseRepository[model.ReadingHistory]{
			DB: db,
		},
	}
}
