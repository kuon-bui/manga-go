package readinghistoryservice

import (
	"manga-go/internal/pkg/logger"
	readinghistoryrepo "manga-go/internal/pkg/repo/readinghistory"

	"go.uber.org/fx"
)

type ReadingHistoryService struct {
	logger             *logger.Logger
	readingHistoryRepo *readinghistoryrepo.ReadingHistoryRepo
}

type ReadingHistoryServiceParams struct {
	fx.In

	Logger             *logger.Logger
	ReadingHistoryRepo *readinghistoryrepo.ReadingHistoryRepo
}

func NewReadingHistoryService(p ReadingHistoryServiceParams) *ReadingHistoryService {
	return &ReadingHistoryService{
		logger:             p.Logger,
		readingHistoryRepo: p.ReadingHistoryRepo,
	}
}
