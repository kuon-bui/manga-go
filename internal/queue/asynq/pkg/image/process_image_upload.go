package image

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"manga-go/internal/pkg/fileprocess"

	"github.com/hibiken/asynq"
)

func (w *ImageProcessWorker) imageProcessHandler(ctx context.Context, task *asynq.Task) error {
	var payload fileprocess.ImageProcessPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		w.logger.Errorf("Failed to unmarshal image process payload: %v", err)
		return err
	}

	if strings.TrimSpace(payload.FilePath) == "" || strings.TrimSpace(payload.TemporaryObjectKey) == "" {
		w.logger.Error("Invalid image process payload")
		return asynq.SkipRetry
	}

	sourceData, err := w.fileService.GetFile(ctx, payload.TemporaryObjectKey)
	if err != nil {
		w.logger.Errorf("Failed to read temporary object %s: %v", payload.TemporaryObjectKey, err)
		if w.fileService.IsNotFoundError(err) || strings.EqualFold(err.Error(), "invalid filename") {
			return asynq.SkipRetry
		}
		return err
	}

	_, err = w.fileService.UploadImageVariants(ctx, payload.FilePath, bytes.NewReader(sourceData))
	if err != nil {
		w.logger.Errorf("Failed to process image upload for %s: %v", payload.FilePath, err)
		return err
	}

	if err := w.fileService.DeleteFile(ctx, payload.TemporaryObjectKey); err != nil {
		w.logger.Warnf("Failed to delete temporary object %s: %v", payload.TemporaryObjectKey, err)
	}

	return nil
}
