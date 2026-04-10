package ratingservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (s *RatingService) ListRatings(ctx context.Context, userID uuid.UUID, paging *common.Paging) response.Result {
	ratings, total, err := s.ratingRepo.FindPaginated(ctx, []any{
		clause.Eq{Column: "user_id", Value: userID},
	}, paging, nil)
	if err != nil {
		s.logger.Error("Failed to list ratings", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultPaginationData(ratings, total, "Ratings retrieved successfully")
}
