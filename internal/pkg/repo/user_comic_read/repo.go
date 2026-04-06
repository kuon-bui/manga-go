package usercomicreadrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type UserComicReadRepo struct {
	*base.BaseRepository[model.UserComicRead]
}

func NewUserComicReadRepo(db *gorm.DB) *UserComicReadRepo {
	return &UserComicReadRepo{
		BaseRepository: &base.BaseRepository[model.UserComicRead]{
			DB: db,
		},
	}
}
