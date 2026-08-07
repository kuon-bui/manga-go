package authorizationadmin

import (
	"context"
	"time"

	authorizationaudit "manga-go/internal/pkg/repo/authorization_audit"

	"github.com/google/uuid"
)

type ListAuditInput struct {
	Page       int
	Limit      int
	Actor      string
	Action     string
	TargetType string
	TargetID   *uuid.UUID
	StartAt    *time.Time
	EndAt      *time.Time
}

func (s *Service) ListAuditLogs(ctx context.Context, input ListAuditInput) (*PagedAuditLogs, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 20
	}
	entries, total, err := s.auditRepo.List(ctx, authorizationaudit.ListInput{
		Page:       input.Page,
		Limit:      input.Limit,
		Actor:      input.Actor,
		Action:     input.Action,
		TargetType: input.TargetType,
		TargetID:   input.TargetID,
		StartAt:    input.StartAt,
		EndAt:      input.EndAt,
	})
	if err != nil {
		return nil, err
	}
	return &PagedAuditLogs{
		Data:  entries,
		Total: total,
		Page:  input.Page,
		Limit: input.Limit,
	}, nil
}
