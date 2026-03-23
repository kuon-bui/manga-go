package comicrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type ComicRepo struct {
	*base.BaseRepository[model.Comic]
}

func NewComicRepo(db *gorm.DB) *ComicRepo {
	return &ComicRepo{
		BaseRepository: &base.BaseRepository[model.Comic]{
			DB: db,
		},
	}
}
