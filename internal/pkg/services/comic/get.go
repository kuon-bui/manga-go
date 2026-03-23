package comicservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ComicService) GetComic(ctx context.Context, slug string) response.Result {
	comic, err := s.comicRepo.FindOne(ctx, []any{
		clause.Eq{Column: "slug", Value: slug},
	}, map[string]common.MoreKeyOption{
		"Authors":        {},
		"Genres":         {},
		"Tags":           {},
		"Chapters.Pages": {},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Comic")
		}
		s.logger.Error("Failed to find comic", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comic retrieved successfully", comic)
}
