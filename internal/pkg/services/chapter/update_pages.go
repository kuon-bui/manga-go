package chapterserivce

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ChapterService) UpdateChapterPages(ctx context.Context, chapterSlug string, req *chapterrequest.UpdateChapterPagesRequest) response.Result {
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

	newPages := make([]*model.Page, len(req.Pages))
	for i, imageURL := range req.Pages {
		newPages[i] = &model.Page{
			ChapterID:  chapter.ID,
			PageNumber: i + 1,
			ImageURL:   imageURL,
		}
	}

	err = s.chapterRepo.UpdateChapterPages(ctx, chapterSlug, newPages)
	if err != nil {
		s.logger.Error("Failed to update chapter pages", "error", err)
		return response.ResultErrDb(err)
	}

	chapter.Pages = newPages

	return response.ResultSuccess("Chapter pages updated successfully", chapter)
}
