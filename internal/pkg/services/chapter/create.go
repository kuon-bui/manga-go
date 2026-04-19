package chapterserivce

import (
	"context"
	"fmt"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	chapterrequest "manga-go/internal/pkg/request/chapter"
	"manga-go/internal/pkg/utils"
	"strings"
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

	currentUser, err := utils.GetCurrentUserFormContext(ctx)
	if err != nil {
		s.logger.Error("Failed to get current user from context", "error", err)
		return response.ResultErrInternal(err)
	}

	slug := req.Slug
	if slug == "" {
		slug = utils.Slugify(req.Title)
	}

	chapter := model.Chapter{
		ComicID:      comicID,
		Number:       req.Number,
		ChapterIdx:   chapterIdx,
		Title:        req.Title,
		Slug:         slug,
		UploadedByID: &currentUser.ID,
	}

	for i, page := range req.Pages {
		pageType := page.PageType
		if pageType == "" {
			pageType = common.ContentTypeImage
		}

		newPage := &model.Page{
			PageNumber: i + 1,
			PageType:   pageType,
		}

		switch pageType {
		case common.ContentTypeImage:
			if strings.TrimSpace(page.ImageURL) == "" {
				return response.ResultError(fmt.Sprintf("page %d: imageUrl is required for image page", i+1))
			}
			newPage.ImageURL = strings.TrimSpace(page.ImageURL)
		case common.ContentTypeText:
			if strings.TrimSpace(page.Content) == "" {
				return response.ResultError(fmt.Sprintf("page %d: content is required for text page", i+1))
			}
			newPage.Content = strings.TrimSpace(page.Content)
		default:
			return response.ResultError(fmt.Sprintf("page %d: invalid pageType", i+1))
		}

		chapter.Pages = append(chapter.Pages, newPage)
	}

	if err := s.chapterRepo.Create(ctx, &chapter); err != nil {
		s.logger.Error("Failed to create chapter", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Chapter created successfully", chapter)
}
