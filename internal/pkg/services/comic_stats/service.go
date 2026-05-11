package comicstatsservice

import (
	"context"
	"errors"
	asynqclient "manga-go/internal/pkg/asynq"
	"manga-go/internal/pkg/logger"
	comicrepo "manga-go/internal/pkg/repo/comic"
	queueconstant "manga-go/internal/queue/queue_constant"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type ComicStatsService struct {
	logger      *logger.Logger
	asynqClient *asynq.Client
	comicRepo   *comicrepo.ComicRepo
}

type ComicStatsServiceParams struct {
	fx.In
	Logger      *logger.Logger
	AsynqClient *asynq.Client
	ComicRepo   *comicrepo.ComicRepo
}

func NewComicStatsService(p ComicStatsServiceParams) *ComicStatsService {
	return &ComicStatsService{
		logger:      p.Logger,
		asynqClient: p.AsynqClient,
		comicRepo:   p.ComicRepo,
	}
}

func (s *ComicStatsService) TriggerUpdateForComic(ctx context.Context, comicID uuid.UUID) error {
	task, err := asynqclient.NewComicStatsUpdateTask(comicID)
	if err != nil {
		s.logger.Error("Failed to create comic stats update task", "error", err)
		return err
	}
	opts := queueconstant.UniqQueue()
	opts = append(opts, asynq.Queue(queueconstant.COMIC_STATS_UPDATE_QUEUE))
	info, err := s.asynqClient.Enqueue(task, opts...)
	if err != nil {
		if errors.Is(err, asynq.ErrDuplicateTask) {
			s.logger.Debugf("Comic stats update task already queued for comic %s", comicID)
			return nil
		}
		s.logger.Error("Failed to enqueue comic stats update task", "error", err)
		return err
	}

	s.logger.Infof("Enqueued comic stats update task for comic %s (task ID: %s)", comicID, info.ID)
	return nil
}

func (s *ComicStatsService) TriggerUpdateForAll(ctx context.Context) (int, error) {
	comics, err := s.comicRepo.FindAll(ctx, nil, nil)
	if err != nil {
		s.logger.Error("Failed to fetch all comics", "error", err)
		return 0, err
	}

	for _, comic := range comics {
		task, err := asynqclient.NewComicStatsUpdateTask(comic.ID)
		if err != nil {
			s.logger.Error("Failed to create comic stats update task", "comicID", comic.ID, "error", err)
			continue
		}

		opts := queueconstant.UniqQueue()
		opts = append(opts, asynq.Queue(queueconstant.COMIC_STATS_UPDATE_QUEUE))
		_, err = s.asynqClient.Enqueue(task, opts...)
		if err != nil {
			if !errors.Is(err, asynq.ErrDuplicateTask) {
				s.logger.Error("Failed to enqueue comic stats update task", "comicID", comic.ID, "error", err)
			}
			continue
		}
	}

	return len(comics), nil
}
