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

func (s *ChapterService) UpdateChapter(ctx context.Context, chapterSlug string, req *chapterrequest.UpdateChapterRequest) response.Result {
	comicID, ok := common.GetComicIdFromContext(ctx)
	if !ok {
		return response.ResultError("Comic not found in context")
	}

	chapter, err := s.chapterRepo.FindOne(ctx, []any{
		clause.Eq{Column: "comic_id", Value: comicID},
		clause.Eq{Column: "slug", Value: chapterSlug},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Chapter")
		}
		s.logger.Error("Failed to find chapter", "error", err)
		return response.ResultErrDb(err)
	}

	updateData := map[string]any{
		"number": req.Number,
		"title":  req.Title,
		"slug":   req.Slug,
	}

	if req.IsPublished != nil {
		updateData["is_published"] = *req.IsPublished
	}

	if err := s.chapterRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: chapter.ID},
	}, updateData); err != nil {
		s.logger.Error("Failed to update chapter", "error", err)
		return response.ResultErrDb(err)
	}

	chapter.Number = req.Number
	chapter.Title = req.Title
	chapter.Slug = req.Slug
	if req.IsPublished != nil {
		chapter.IsPublished = *req.IsPublished
	}

	return response.ResultSuccess("Chapter updated successfully", chapter)
}
