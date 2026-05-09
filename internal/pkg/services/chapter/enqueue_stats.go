package chapterserivce

import (
	"time"

	asynqclient "manga-go/internal/pkg/asynq"
	queueconstant "manga-go/internal/queue/queue_constant"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

func (s *ChapterService) enqueueStatsUpdate(comicID uuid.UUID) {
	if s.asynqClient == nil {
		return
	}

	task, err := asynqclient.NewComicStatsUpdateTask(comicID)
	if err != nil {
		s.logger.Error("Failed to create comic stats update task", "error", err)
		return
	}

	if _, err := s.asynqClient.Enqueue(task,
		asynq.Queue(queueconstant.COMIC_STATS_UPDATE_QUEUE),
		asynq.Unique(30*time.Second),
	); err != nil {
		s.logger.Error("Failed to enqueue comic stats update", "comicID", comicID, "error", err)
	}
}
