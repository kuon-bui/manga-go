package authorizationaudit

import (
	"context"
	"strings"
	"time"

	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ListInput struct {
	Page       int
	Limit      int
	Actor      string
	Action     string
	TargetType string
	TargetID   *uuid.UUID
	StartAt    *time.Time
	EndAt      *time.Time
}

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) AppendTx(tx *gorm.DB, entry *model.AuthorizationAuditLog) error {
	return tx.Create(entry).Error
}

func (r *Repo) List(ctx context.Context, input ListInput) ([]*model.AuthorizationAuditLog, int64, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 20
	}
	query := r.db.WithContext(ctx).Model(&model.AuthorizationAuditLog{})
	if actor := strings.TrimSpace(strings.ToLower(input.Actor)); actor != "" {
		pattern := "%" + actor + "%"
		query = query.Where(
			"LOWER(actor_name_snapshot) LIKE ? OR LOWER(actor_email_snapshot) LIKE ?",
			pattern,
			pattern,
		)
	}
	if input.Action != "" {
		query = query.Where("action = ?", input.Action)
	}
	if input.TargetType != "" {
		query = query.Where("target_type = ?", input.TargetType)
	}
	if input.TargetID != nil {
		query = query.Where("target_id = ?", *input.TargetID)
	}
	if input.StartAt != nil {
		query = query.Where("created_at >= ?", input.StartAt.UTC())
	}
	if input.EndAt != nil {
		query = query.Where("created_at <= ?", input.EndAt.UTC())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var entries []*model.AuthorizationAuditLog
	if err := query.
		Order("created_at DESC, id DESC").
		Limit(input.Limit).
		Offset((input.Page - 1) * input.Limit).
		Find(&entries).Error; err != nil {
		return nil, 0, err
	}
	return entries, total, nil
}
