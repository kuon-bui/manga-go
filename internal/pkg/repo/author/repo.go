package authorrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type AuthorRepo struct {
	*base.BaseRepository[model.Author]
}

func NewAuthorRepo(db *gorm.DB) *AuthorRepo {
	return &AuthorRepo{
		BaseRepository: &base.BaseRepository[model.Author]{
			DB: db,
		},
	}
}
