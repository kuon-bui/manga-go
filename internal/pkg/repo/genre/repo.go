package genrerepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type GenreRepo struct {
	*base.BaseRepository[model.Genre]
}

func NewGenreRepo(db *gorm.DB) *GenreRepo {
	return &GenreRepo{
		BaseRepository: &base.BaseRepository[model.Genre]{
			DB: db,
		},
	}
}
