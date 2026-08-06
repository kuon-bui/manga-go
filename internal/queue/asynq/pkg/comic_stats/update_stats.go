package comicstatsworker

import (
	"context"
	"encoding/json"

	asynqclient "manga-go/internal/pkg/asynq"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

func (w *ComicStatsWorker) updateStatsHandler(ctx context.Context, task *asynq.Task) error {
	var payload asynqclient.ComicStatsUpdatePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		w.logger.Errorf("Failed to unmarshal comic stats update payload: %v", err)
		return err
	}

	if payload.ComicID == uuid.Nil {
		w.logger.Error("Invalid comic stats update payload: missing comicID")
		return asynq.SkipRetry
	}

	if err := w.statsService.RecomputeComicStats(ctx, payload.ComicID); err != nil {
		w.logger.Errorf("Failed to recompute comic stats for %s: %v", payload.ComicID, err)
		return err
	}

	w.logger.Infof("Comic stats updated successfully for comic %s", payload.ComicID)
	return nil
}
