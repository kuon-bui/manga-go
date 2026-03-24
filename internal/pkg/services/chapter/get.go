package chapterserivce

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ChapterService) GetChapter(ctx context.Context, chapterSlug string) response.Result {
	comicID, ok := common.GetComicIdFromContext(ctx)
	if !ok {
		return response.ResultError("Comic not found in context")
	}

	chapter, err := s.chapterRepo.FindOne(ctx, []any{
		clause.Eq{Column: "comic_id", Value: comicID},
		clause.Eq{Column: "slug", Value: chapterSlug},
	}, map[string]common.MoreKeyOption{
		"Pages": {},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Chapter")
		}
		s.logger.Error("Failed to find chapter", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Chapter retrieved successfully", chapter)
}
