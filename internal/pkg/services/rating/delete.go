package ratingservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *RatingService) DeleteRating(ctx context.Context, userID, id uuid.UUID) response.Result {
	_, err := s.ratingRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
		clause.Eq{Column: "user_id", Value: userID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Rating")
		}
		s.logger.Error("Failed to find rating for deletion", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.ratingRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "id", Value: id},
		clause.Eq{Column: "user_id", Value: userID},
	}); err != nil {
		s.logger.Error("Failed to delete rating", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Rating deleted successfully", nil)
}
