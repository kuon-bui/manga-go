package readinghistoryservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	readinghistoryrequest "manga-go/internal/pkg/request/reading_history"

	"time"

	"github.com/google/uuid"
)

func (s *ReadingHistoryService) CreateReadingHistory(ctx context.Context, userID uuid.UUID, req *readinghistoryrequest.CreateReadingHistoryRequest) response.Result {
	now := time.Now()
	readingHistory := model.ReadingHistory{
		UserID:     userID,
		ChapterID:  req.ChapterID,
		ComicID:    req.ComicID,
		LastReadAt: &now,
	}

	if err := s.readingHistoryRepo.Create(ctx, &readingHistory); err != nil {
		s.logger.Error("Failed to create reading history", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Reading history created successfully", readingHistory)
}
