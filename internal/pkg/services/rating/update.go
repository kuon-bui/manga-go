package ratingservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	ratingrequest "manga-go/internal/pkg/request/rating"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *RatingService) UpdateRating(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *ratingrequest.UpdateRatingRequest) response.Result {
	rating, err := s.ratingRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
		clause.Eq{Column: "user_id", Value: userID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Rating")
		}
		s.logger.Error("Failed to find rating", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.ratingRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: id},
		clause.Eq{Column: "user_id", Value: userID},
	}, map[string]any{
		"score": req.Score,
	}); err != nil {
		s.logger.Error("Failed to update rating", "error", err)
		return response.ResultErrDb(err)
	}

	rating.Score = req.Score
	return response.ResultSuccess("Rating updated successfully", rating)
}
