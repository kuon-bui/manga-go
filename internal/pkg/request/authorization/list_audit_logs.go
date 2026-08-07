package authorizationrequest

import (
	"errors"
	"time"

	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type ListAuditLogsRequest struct {
	common.Paging
	Actor      string     `form:"actor" binding:"omitempty,max=255"`
	Action     string     `form:"action" binding:"omitempty,max=64"`
	TargetType string     `form:"target_type" binding:"omitempty,max=32"`
	TargetID   *uuid.UUID `form:"target_id"`
	StartAt    *time.Time `form:"start_at" time_format:"2006-01-02T15:04:05Z07:00"`
	EndAt      *time.Time `form:"end_at" time_format:"2006-01-02T15:04:05Z07:00"`
}

func (r *ListAuditLogsRequest) Validate() error {
	r.Fulfill()
	if r.Limit > 100 {
		return errors.New("limit must not exceed 100")
	}
	if r.StartAt != nil && r.EndAt != nil && r.StartAt.After(*r.EndAt) {
		return errors.New("start_at must not be after end_at")
	}
	return nil
}
