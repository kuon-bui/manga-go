package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ComicService) UnfollowComic(ctx context.Context, userID, comicID uuid.UUID) response.Result {
	_, err := s.comicFollowRepo.FindByUserAndComic(ctx, userID, comicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("ComicFollow")
		}

		s.logger.Error("Failed to find comic follow for deletion", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.comicFollowRepo.DeletePermanently(ctx, []any{
		clause.Eq{Column: "user_id", Value: userID},
		clause.Eq{Column: "comic_id", Value: comicID},
	}); err != nil {
		s.logger.Error("Failed to delete comic follow", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comic unfollowed successfully", nil)
}
