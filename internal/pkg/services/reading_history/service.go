package readinghistoryservice

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	readinghistoryrepo "manga-go/internal/pkg/repo/reading_history"

	"go.uber.org/fx"
)

// ReadingHistoryRepository defines the data access interface for ReadingHistory.
type ReadingHistoryRepository interface {
	Create(ctx context.Context, readingHistory *model.ReadingHistory) error
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.ReadingHistory, error)
	Update(ctx context.Context, conditions []any, data map[string]any) error
	DeleteSoft(ctx context.Context, conditions []any) error
	FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.ReadingHistory, int64, error)
}

type ReadingHistoryService struct {
	logger             *logger.Logger
	readingHistoryRepo ReadingHistoryRepository
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

// NewReadingHistoryServiceWithRepo creates a ReadingHistoryService with an explicit repository,
// useful for unit testing.
func NewReadingHistoryServiceWithRepo(l *logger.Logger, repo ReadingHistoryRepository) *ReadingHistoryService {
	return &ReadingHistoryService{
		logger:             l,
		readingHistoryRepo: repo,
	}
}
