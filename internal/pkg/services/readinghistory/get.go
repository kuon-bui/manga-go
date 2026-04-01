package readinghistoryservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ReadingHistoryService) GetReadingHistory(ctx context.Context, id uuid.UUID) response.Result {
	readingHistory, err := s.readingHistoryRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("ReadingHistory")
		}
		s.logger.Error("Failed to get reading history", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Reading history retrieved successfully", readingHistory)
}
