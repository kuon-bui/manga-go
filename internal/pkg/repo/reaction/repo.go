package reactionrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type ReactionRepo struct {
	*base.BaseRepository[model.Reaction]
}

func NewReactionRepo(db *gorm.DB) *ReactionRepo {
	return &ReactionRepo{
		BaseRepository: &base.BaseRepository[model.Reaction]{
			DB: db,
		},
	}
}
