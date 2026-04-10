package ratingrepo

import (
	"context"
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RatingRepo struct {
	*base.BaseRepository[model.Rating]
}

func NewRatingRepo(db *gorm.DB) *RatingRepo {
	return &RatingRepo{
		BaseRepository: &base.BaseRepository[model.Rating]{
			DB: db,
		},
	}
}

func (r *RatingRepo) FindByUserAndComic(ctx context.Context, userID uuid.UUID, comicID uuid.UUID) (*model.Rating, error) {
	return r.FindOne(ctx, []any{
		clause.Eq{Column: "user_id", Value: userID},
		clause.Eq{Column: "comic_id", Value: comicID},
	}, nil)
}
