package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ComicService) DeleteComic(ctx context.Context, slug string) response.Result {
	comic, err := s.comicRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Comic")
		}
		s.logger.Error("Failed to find comic for deletion", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.comicRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "id", Value: comic.ID},
	}); err != nil {
		s.logger.Error("Failed to delete comic", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comic deleted successfully", nil)
}
