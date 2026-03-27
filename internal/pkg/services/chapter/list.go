package chapterserivce

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"gorm.io/gorm/clause"
)

func (s *ChapterService) ListChapters(ctx context.Context, paging *common.Paging) response.Result {
	comicID, ok := common.GetComicIdFromContext(ctx)
	if !ok {
		return response.ResultError("Comic not found in context")
	}

	chapters, total, err := s.chapterRepo.FindPaginated(ctx, []any{
		clause.Eq{Column: "comic_id", Value: comicID},
	}, paging, nil)
	if err != nil {
		s.logger.Error("Failed to list chapters", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Chapters retrieved successfully", response.ResponsePaginationData(chapters, total))
}
