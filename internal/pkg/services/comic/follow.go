package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/constant"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ComicService) FollowComic(ctx context.Context, userID, comicID uuid.UUID, followStatus constant.FollowStatus) response.Result {
	if followStatus == "" {
		followStatus = constant.FollowStatusReading
	}

	existingFollow, err := s.comicFollowRepo.FindByUserAndComicWithUnscoped(ctx, userID, comicID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to find comic follow", "error", err)
		return response.ResultErrDb(err)
	}

	if existingFollow != nil {
		if existingFollow.DeletedAt.Valid {
			if err := s.comicFollowRepo.RestoreByUserAndComic(ctx, userID, comicID); err != nil {
				s.logger.Error("Failed to restore comic follow", "error", err)
				return response.ResultErrDb(err)
			}

			if err := s.comicFollowRepo.Update(ctx, []any{
				clause.Eq{Column: "user_id", Value: userID},
				clause.Eq{Column: "comic_id", Value: comicID},
			}, map[string]any{
				"follow_status": followStatus,
			}); err != nil {
				s.logger.Error("Failed to update follow status after restore", "error", err)
				return response.ResultErrDb(err)
			}

			restoredFollow, err := s.comicFollowRepo.FindByUserAndComic(ctx, userID, comicID)
			if err != nil {
				s.logger.Error("Failed to reload restored comic follow", "error", err)
				return response.ResultErrDb(err)
			}

			return response.ResultSuccess("Comic followed successfully", restoredFollow)
		}

		if existingFollow.FollowStatus != followStatus {
			if err := s.comicFollowRepo.Update(ctx, []any{
				clause.Eq{Column: "user_id", Value: userID},
				clause.Eq{Column: "comic_id", Value: comicID},
			}, map[string]any{
				"follow_status": followStatus,
			}); err != nil {
				s.logger.Error("Failed to update comic follow status", "error", err)
				return response.ResultErrDb(err)
			}

			existingFollow.FollowStatus = followStatus
			return response.ResultSuccess("Comic follow status updated successfully", existingFollow)
		}

		return response.ResultSuccess("Comic already followed", existingFollow)
	}

	comicFollow := model.ComicFollow{
		UserID:       userID,
		ComicID:      comicID,
		FollowStatus: followStatus,
	}

	if err := s.comicFollowRepo.Create(ctx, &comicFollow); err != nil {
		s.logger.Error("Failed to create comic follow", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comic followed successfully", comicFollow)
}
