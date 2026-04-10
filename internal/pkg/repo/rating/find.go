package ratingrepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

func (r *RatingRepo) GetAverageRatingByComicID(ctx context.Context, comicID uuid.UUID) (float64, int64, error) {
	var result struct {
		Average float64
		Count   int64
	}

	err := r.DB.WithContext(ctx).Model(&model.Rating{}).
		Select("AVG(score) as average, COUNT(*) as count").
		Where("comic_id = ?", comicID).
		Scan(&result).Error

	if err != nil {
		return 0, 0, err
	}

	return result.Average, result.Count, nil
}
