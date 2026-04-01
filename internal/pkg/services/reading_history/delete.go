package readinghistoryservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ReadingHistoryService) DeleteReadingHistory(ctx context.Context, id uuid.UUID) response.Result {
	readingHistory, err := s.readingHistoryRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("ReadingHistory")
		}
		s.logger.Error("Failed to find reading history", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.readingHistoryRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}); err != nil {
		s.logger.Error("Failed to delete reading history", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Reading history deleted successfully", readingHistory)
}
