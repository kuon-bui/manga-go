package image

import (
	"context"
	"encoding/json"
	"strings"

	"manga-go/internal/pkg/fileprocess"

	"github.com/hibiken/asynq"
)

func (w *ImageProcessWorker) imageProcessCleanupHandler(ctx context.Context, task *asynq.Task) error {
	var payload fileprocess.ImageProcessCleanupPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		w.logger.Errorf("Failed to unmarshal image process cleanup payload: %v", err)
		return asynq.SkipRetry
	}

	if strings.TrimSpace(payload.TemporaryObjectKey) == "" {
		w.logger.Error("Invalid image process cleanup payload")
		return asynq.SkipRetry
	}

	err := w.fileService.DeleteFile(ctx, payload.TemporaryObjectKey)
	if err == nil {
		return nil
	}

	if w.fileService.IsNotFoundError(err) || strings.EqualFold(err.Error(), "invalid filename") {
		return nil
	}

	w.logger.Errorf("Failed to cleanup temporary object %s: %v", payload.TemporaryObjectKey, err)
	return err
}
