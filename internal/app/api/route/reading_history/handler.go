package readinghistoryroute

import (
	readinghistoryservice "manga-go/internal/pkg/services/reading_history"

	"go.uber.org/fx"
)

type ReadingHistoryHandler struct {
	readingHistoryService *readinghistoryservice.ReadingHistoryService
}

type ReadingHistoryHandlerParams struct {
	fx.In

	ReadingHistoryService *readinghistoryservice.ReadingHistoryService
}

func NewReadingHistoryHandler(p ReadingHistoryHandlerParams) *ReadingHistoryHandler {
	return &ReadingHistoryHandler{
		readingHistoryService: p.ReadingHistoryService,
	}
}
