package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComicFollowStatus struct {
	IsFollowed bool `json:"isFollowed"`
}

func (s *ComicService) GetFollowStatus(ctx context.Context, userID, comicID uuid.UUID) response.Result {
	_, err := s.comicFollowRepo.FindByUserAndComic(ctx, userID, comicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultSuccess("Comic follow status retrieved successfully", ComicFollowStatus{IsFollowed: false})
		}

		s.logger.Error("Failed to get comic follow status", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comic follow status retrieved successfully", ComicFollowStatus{IsFollowed: true})
}
