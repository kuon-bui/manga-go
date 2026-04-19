package ratingservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRatingResponse struct {
	Rating  *int    `json:"rating"`
	Comment *string `json:"comment,omitempty"`
}

func (s *RatingService) GetUserRatingForComic(ctx context.Context, userID, comicID uuid.UUID) response.Result {
	rating, err := s.ratingRepo.FindByUserAndComic(ctx, userID, comicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Frontend expects rating: null if not rated
			return response.ResultSuccess("No rating found", UserRatingResponse{Rating: nil})
		}
		s.logger.Error("Failed to find user rating", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Rating retrieved successfully", UserRatingResponse{
		Rating:  &rating.Score,
		Comment: rating.Comment,
	})
}
