package ratingservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	ratingrequest "manga-go/internal/pkg/request/rating"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *RatingService) CreateRating(ctx context.Context, userID, comicId uuid.UUID, req *ratingrequest.CreateRatingRequest) response.Result {
	existingRating, err := s.ratingRepo.FindByUserAndComic(ctx, userID, comicId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to find user rating", "error", err)
		return response.ResultErrDb(err)
	}

	if existingRating != nil {
		if err := s.ratingRepo.Update(ctx, []any{
			clause.Eq{Column: "id", Value: existingRating.ID},
			clause.Eq{Column: "user_id", Value: userID},
		}, map[string]any{
			"score": req.Score,
		}); err != nil {
			s.logger.Error("Failed to update rating", "error", err)
			return response.ResultErrDb(err)
		}

		existingRating.Score = req.Score
		return response.ResultSuccess("Rating updated successfully", existingRating)
	}

	rating := model.Rating{
		UserId:  userID,
		ComicId: comicId,
		Score:   req.Score,
	}

	if err := s.ratingRepo.Create(ctx, &rating); err != nil {
		s.logger.Error("Failed to create rating", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Rating created successfully", rating)
}
