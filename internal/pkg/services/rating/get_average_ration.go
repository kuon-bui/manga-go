package ratingservice

import (
	"context"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
)

func (s *RatingService) GetAverageRating(ctx context.Context, comicId uuid.UUID) response.Result {
	average, count, err := s.ratingRepo.GetAverageRatingByComicID(ctx, comicId)
	if err != nil {
		s.logger.Error("Failed to get average rating", "error", err)
		return response.ResultErrDb(err)
	}

	totalRating := map[string]any{
		"average": average,
		"count":   count,
	}

	return response.ResultSuccess("Average rating retrieved successfully", totalRating)
}
