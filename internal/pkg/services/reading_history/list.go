package readinghistoryservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (s *ReadingHistoryService) ListReadingHistories(ctx context.Context, userID uuid.UUID, paging *common.Paging) response.Result {
	readingHistories, total, err := s.readingHistoryRepo.FindPaginated(ctx, []any{
		clause.Eq{Column: "user_id", Value: userID},
	}, paging, nil)
	if err != nil {
		s.logger.Error("Failed to list reading histories", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Reading histories retrieved successfully", response.ResponsePaginationData(readingHistories, total))
}
