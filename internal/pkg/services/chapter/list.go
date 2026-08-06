package chapterserivce

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/common"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ChapterService) ListChapters(ctx context.Context, paging *common.Paging) response.Result {
	comicID, ok := common.GetComicIdFromContext(ctx)
	if !ok {
		return response.ResultError("Comic not found in context")
	}

	user := authorization.ViewerFromContext(ctx).User
	conditions := []any{
		clause.Eq{Column: "comic_id", Value: comicID},
	}
	if user == nil {
		conditions = append(conditions, s.chapterRepo.IsPublished(true))
	}

	chapters, total, err := s.chapterRepo.FindPaginated(ctx, conditions, paging, nil)
	if err != nil {
		s.logger.Error("Failed to list chapters", "error", err)
		return response.ResultErrDb(err)
	}
	if user == nil {
		return response.ResultPaginationData(chapters, total, "Chapters retrieved successfully")
	}

	userComicRead, err := s.userComicReadRepo.FindOne(ctx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "comic_id", Value: comicID},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to find user comic read data", "error", err)
		return response.ResultErrDb(err)
	}

	if userComicRead != nil {
		for _, chapter := range chapters {
			chapter.IsRead = userComicRead.ReadData.IsRead(int(chapter.ChapterIdx))
		}
	}

	return response.ResultPaginationData(chapters, total, "Chapters retrieved successfully")
}
