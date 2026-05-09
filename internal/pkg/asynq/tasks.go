package asynqclient

import (
	"encoding/json"

	queueconstant "manga-go/internal/queue/queue_constant"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type ComicStatsUpdatePayload struct {
	ComicID uuid.UUID `json:"comicId"`
}

func NewComicStatsUpdateTask(comicID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(ComicStatsUpdatePayload{ComicID: comicID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(queueconstant.COMIC_STATS_UPDATE_TASK, payload), nil
}
