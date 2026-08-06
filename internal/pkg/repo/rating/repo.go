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

func (r *RatingRepo) FindAvgRatingByComicID(ctx context.Context, comicID uuid.UUID) (*float64, error) {
	return r.findAvgRating(r.DB.WithContext(ctx), comicID)
}

func (r *RatingRepo) FindAvgRatingByComicIDWithTransaction(tx *gorm.DB, comicID uuid.UUID) (*float64, error) {
	return r.findAvgRating(tx, comicID)
}

func (r *RatingRepo) findAvgRating(db *gorm.DB, comicID uuid.UUID) (*float64, error) {
	type result struct {
		Avg *float64 `gorm:"column:avg"`
	}
	var res result
	if err := db.
		Model(&model.Rating{}).
		Select("AVG(score) as avg").
		Where("comic_id = ?", comicID).
		Scan(&res).Error; err != nil {
		return nil, err
	}
	return res.Avg, nil
}
