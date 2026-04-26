package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/constant"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ComicService) UpdateComicFollowStatus(ctx context.Context, userID, comicID uuid.UUID, followStatus constant.FollowStatus) response.Result {
	follow, err := s.comicFollowRepo.FindByUserAndComic(ctx, userID, comicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("ComicFollow")
		}

		s.logger.Error("Failed to find comic follow for status update", "error", err)
		return response.ResultErrDb(err)
	}

	if follow.FollowStatus == followStatus {
		return response.ResultSuccess("Comic follow status updated successfully", follow)
	}

	if err := s.comicFollowRepo.Update(ctx, []any{
		clause.Eq{Column: "user_id", Value: userID},
		clause.Eq{Column: "comic_id", Value: comicID},
	}, map[string]any{
		"follow_status": followStatus,
	}); err != nil {
		s.logger.Error("Failed to update comic follow status", "error", err)
		return response.ResultErrDb(err)
	}

	follow.FollowStatus = followStatus
	return response.ResultSuccess("Comic follow status updated successfully", follow)
}
