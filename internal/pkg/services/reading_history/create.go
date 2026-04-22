package readinghistoryservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	readinghistoryrequest "manga-go/internal/pkg/request/reading_history"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ReadingHistoryService) CreateReadingHistory(ctx context.Context, userID uuid.UUID, req *readinghistoryrequest.CreateReadingHistoryRequest) response.Result {
	now := time.Now()
	conditions := []any{
		clause.Eq{Column: "user_id", Value: userID},
		clause.Eq{Column: "chapter_id", Value: req.ChapterID},
	}

	existingReadingHistory, err := s.readingHistoryRepo.FindOne(ctx, conditions, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to find reading history", "error", err)
		return response.ResultErrDb(err)
	}

	if existingReadingHistory != nil && existingReadingHistory.ID != uuid.Nil {
		existingReadingHistory.LastReadAt = &now
		if err := s.readingHistoryRepo.Update(ctx, conditions, map[string]any{
			"last_read_at": now,
		}); err != nil {
			s.logger.Error("Failed to update reading history", "error", err)
			return response.ResultErrDb(err)
		}

		return response.ResultSuccess("Reading history updated successfully", existingReadingHistory)
	}

	readingHistory := model.ReadingHistory{
		UserID:     userID,
		ChapterID:  req.ChapterID,
		ComicID:    req.ComicID,
		LastReadAt: &now,
	}

	if err := s.readingHistoryRepo.Save(ctx, &readingHistory); err != nil {
		s.logger.Error("Failed to save reading history", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Reading history created successfully", readingHistory)
}
