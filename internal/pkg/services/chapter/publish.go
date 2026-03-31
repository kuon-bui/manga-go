package chapterserivce

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ChapterService) PublishChapter(ctx context.Context, chapterSlug string, req *chapterrequest.PublishChapterRequest) response.Result {
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

	if req.IsPublished && len(chapter.Pages) == 0 {
		return response.ResultError("Chapter must have at least one page before publishing")
	}

	if err := s.chapterRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: chapter.ID},
	}, map[string]any{
		"is_published": req.IsPublished,
	}); err != nil {
		s.logger.Error("Failed to publish chapter", "error", err)
		return response.ResultErrDb(err)
	}

	chapter.IsPublished = req.IsPublished

	msg := "Chapter unpublished successfully"
	if req.IsPublished {
		msg = "Chapter published successfully"
	}

	return response.ResultSuccess(msg, chapter)
}
