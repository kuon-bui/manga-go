package notification

import (
	"context"
	"encoding/json"
	notificationpkg "manga-go/internal/pkg/notification"

	"github.com/hibiken/asynq"
)

func (w *NotificationWorker) notificationFanoutHandler(ctx context.Context, task *asynq.Task) error {
	var payload notificationpkg.FanoutPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		w.logger.Errorf("Failed to unmarshal notification fanout payload: %v", err)
		return err
	}

	return w.notificationService.HandleFanout(ctx, &payload)
}
