package pagereactionrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type PageReactionRepo struct {
	*base.BaseRepository[model.PageReaction]
}

func NewPageReactionRepo(db *gorm.DB) *PageReactionRepo {
	return &PageReactionRepo{
		BaseRepository: &base.BaseRepository[model.PageReaction]{
			DB: db,
		},
	}
}
