package pagerepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type PageRepo struct {
	*base.BaseRepository[model.Page]
}

func NewPageRepo(db *gorm.DB) *PageRepo {
	return &PageRepo{
		BaseRepository: &base.BaseRepository[model.Page]{
			DB: db,
		},
	}
}
