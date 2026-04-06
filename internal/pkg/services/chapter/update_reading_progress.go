package chapterserivce

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ChapterService) UpdateReadingProgress(
	ctx context.Context,
	user *model.User,
	chapterID uuid.UUID,
	req *chapterrequest.UpdateReadingProgressRequest,
) response.Result {
	chapter, err := s.chapterRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: chapterID},
	}, nil)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ResultNotFound("Chapter")
		}
		s.logger.Error("Failed to find chapter", "error", err)
		return response.ResultErrDb(err)
	}

	readingProgress, err := s.readingProgressRepo.FindOne(ctx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "comic_id", Value: chapter.ComicID},
	}, nil)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			readingProgress = &model.ReadingProgress{
				UserID:    user.ID,
				ComicID:   chapter.ComicID,
				ChapterID: chapter.ID,
			}
		} else {
			s.logger.Error("Failed to find reading progress", "error", err)
			return response.ResultErrDb(err)
		}
	}

	readingProgress.ScrollPercent = req.ScrollPercent

	err = s.readingProgressRepo.Save(ctx, readingProgress)
	if err != nil {
		s.logger.Error("Failed to save reading progress", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Reading progress updated successfully", nil)
}
