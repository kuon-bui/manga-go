package notificationrepo

import (
	"context"
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NotificationRepo struct {
	*base.BaseRepository[model.Notification]
}

func NewNotificationRepo(db *gorm.DB) *NotificationRepo {
	return &NotificationRepo{
		BaseRepository: &base.BaseRepository[model.Notification]{
			DB: db,
		},
	}
}

func (r *NotificationRepo) FindByDedupeKey(ctx context.Context, dedupeKey string) (*model.Notification, error) {
	return r.FindOne(ctx, []any{
		clause.Eq{Column: "dedupe_key", Value: dedupeKey},
	}, nil)
}

func (r *NotificationRepo) FindByDedupeKeyWithTransaction(tx *gorm.DB, dedupeKey string) (*model.Notification, error) {
	return r.FindOneWithTransaction(tx, []any{
		clause.Eq{Column: "dedupe_key", Value: dedupeKey},
	}, nil)
}
