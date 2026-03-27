package tagrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type TagRepo struct {
	*base.BaseRepository[model.Tag]
}

func NewTagRepo(db *gorm.DB) *TagRepo {
	return &TagRepo{
		BaseRepository: &base.BaseRepository[model.Tag]{
			DB: db,
		},
	}
}
