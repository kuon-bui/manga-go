package comicstatsworker

import (
	"manga-go/internal/pkg/logger"
	queueconstant "manga-go/internal/queue/queue_constant"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type ComicStatsWorkerParams struct {
	fx.In
	Mux          *asynq.ServeMux
	Logger       *logger.Logger
	StatsService *StatsService
}

type ComicStatsWorker struct {
	mux          *asynq.ServeMux
	logger       *logger.Logger
	statsService *StatsService
}

func NewComicStatsWorker(p ComicStatsWorkerParams) *ComicStatsWorker {
	return &ComicStatsWorker{
		mux:          p.Mux,
		logger:       p.Logger,
		statsService: p.StatsService,
	}
}

func (w *ComicStatsWorker) RegisterWorkers() {
	w.mux.HandleFunc(queueconstant.COMIC_STATS_UPDATE_TASK, w.updateStatsHandler)
}
