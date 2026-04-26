package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/constant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComicFollowStatus struct {
	IsFollowed   bool                   `json:"isFollowed"`
	FollowStatus *constant.FollowStatus `json:"followStatus,omitempty"`
}

func (s *ComicService) GetFollowStatus(ctx context.Context, userID, comicID uuid.UUID) response.Result {
	follow, err := s.comicFollowRepo.FindByUserAndComic(ctx, userID, comicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultSuccess("Comic follow status retrieved successfully", ComicFollowStatus{IsFollowed: false})
		}

		s.logger.Error("Failed to get comic follow status", "error", err)
		return response.ResultErrDb(err)
	}

	followStatus := follow.FollowStatus
	return response.ResultSuccess("Comic follow status retrieved successfully", ComicFollowStatus{IsFollowed: true, FollowStatus: &followStatus})
}
