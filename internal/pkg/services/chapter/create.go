package chapterserivce

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	chapterrequest "manga-go/internal/pkg/request/chapter"
)

func (s *ChapterService) CreateChapter(ctx context.Context, req *chapterrequest.CreateChapterRequest) response.Result {
	comicID, ok := common.GetComicIdFromContext(ctx)
	if !ok {
		return response.ResultError("Comic not found in context")
	}

	chapterIdx, err := s.chapterRepo.GetNextChapterIdx(ctx, comicID)
	if err != nil {
		s.logger.Error("Failed to get next chapter index", "error", err)
		return response.ResultErrDb(err)
	}

	chapter := model.Chapter{
		ComicID:    comicID,
		Number:     req.Number,
		ChapterIdx: chapterIdx,
		Title:      req.Title,
		Slug:       req.Slug,
	}

	for i, page := range req.Pages {
		chapter.Pages = append(chapter.Pages, &model.Page{
			PageNumber: i + 1,
			ImageURL:   page,
		})
	}

	if err := s.chapterRepo.Create(ctx, &chapter); err != nil {
		s.logger.Error("Failed to create chapter", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Chapter created successfully", chapter)
}
