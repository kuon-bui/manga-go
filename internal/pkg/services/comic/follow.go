package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *ComicService) FollowComic(ctx context.Context, userID, comicID uuid.UUID) response.Result {
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

			restoredFollow, err := s.comicFollowRepo.FindByUserAndComic(ctx, userID, comicID)
			if err != nil {
				s.logger.Error("Failed to reload restored comic follow", "error", err)
				return response.ResultErrDb(err)
			}

			return response.ResultSuccess("Comic followed successfully", restoredFollow)
		}

		return response.ResultSuccess("Comic already followed", existingFollow)
	}

	comicFollow := model.ComicFollow{
		UserID:  userID,
		ComicID: comicID,
	}

	if err := s.comicFollowRepo.Create(ctx, &comicFollow); err != nil {
		s.logger.Error("Failed to create comic follow", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comic followed successfully", comicFollow)
}
