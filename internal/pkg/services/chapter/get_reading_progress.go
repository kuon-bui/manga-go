package chapterserivce

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ChapterService) GetReadingProgress(ctx context.Context, user *model.User) response.Result {
	comicID, ok := common.GetComicIdFromContext(ctx)
	if !ok {
		return response.ResultError("Comic not found in context")
	}

	readingProgress, err := s.readingProgressRepo.FindOne(ctx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "comic_id", Value: comicID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Reading progress")
		}

		s.logger.Error("Failed to find reading progress", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Reading progress retrieved successfully", readingProgress)
}
