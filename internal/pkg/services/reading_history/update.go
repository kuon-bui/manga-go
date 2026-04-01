package readinghistoryservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	readinghistoryrequest "manga-go/internal/pkg/request/reading_history"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ReadingHistoryService) UpdateReadingHistory(ctx context.Context, id uuid.UUID, req *readinghistoryrequest.UpdateReadingHistoryRequest) response.Result {
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

	updateData := map[string]any{}
	if req.LastReadAt != nil {
		updateData["last_read_at"] = req.LastReadAt
		readingHistory.LastReadAt = req.LastReadAt
	} else {
		now := time.Now()
		updateData["last_read_at"] = now
		readingHistory.LastReadAt = &now
	}

	if err := s.readingHistoryRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, updateData); err != nil {
		s.logger.Error("Failed to update reading history", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Reading history updated successfully", readingHistory)
}
