package usernotificationrepo

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	pknotification "manga-go/internal/pkg/notification"
	"manga-go/internal/pkg/repo/base"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserNotificationRepo struct {
	*base.BaseRepository[model.UserNotification]
}

func NewUserNotificationRepo(db *gorm.DB) *UserNotificationRepo {
	return &UserNotificationRepo{
		BaseRepository: &base.BaseRepository[model.UserNotification]{
			DB: db,
		},
	}
}

func (r *UserNotificationRepo) FindByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (*model.UserNotification, error) {
	return r.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
		clause.Eq{Column: "user_id", Value: userID},
	}, map[string]common.MoreKeyOption{
		"Notification": {},
	})
}

func (r *UserNotificationRepo) FindPaginatedByUserID(ctx context.Context, userID uuid.UUID, paging *common.Paging, unreadOnly bool, notificationType pknotification.Type) ([]*model.UserNotification, int64, error) {
	var items []*model.UserNotification
	var total int64

	db := r.DB.WithContext(ctx).
		Model(&model.UserNotification{}).
		Joins("JOIN notifications ON notifications.id = user_notifications.notification_id").
		Where("user_notifications.user_id = ?", userID).
		Where("notifications.deleted_at IS NULL").
		Preload("Notification").
		Order("user_notifications.created_at DESC")

	if unreadOnly {
		db = db.Where("user_notifications.is_read = ?", false)
	}

	if notificationType != "" {
		db = db.Where("notifications.type = ?", notificationType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Scopes(r.WithPaginate(paging))

	if err := db.Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *UserNotificationRepo) CreateListIgnoreConflictsWithTransaction(tx *gorm.DB, items []*model.UserNotification) error {
	if len(items) == 0 {
		return nil
	}

	return tx.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(items).Error
}

func (r *UserNotificationRepo) FindByNotificationAndUserIDs(ctx context.Context, notificationID uuid.UUID, userIDs []uuid.UUID) ([]*model.UserNotification, error) {
	if len(userIDs) == 0 {
		return []*model.UserNotification{}, nil
	}

	var items []*model.UserNotification
	err := r.DB.WithContext(ctx).
		Model(&model.UserNotification{}).
		Where("notification_id = ?", notificationID).
		Where("user_id IN ?", userIDs).
		Preload("Notification").
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *UserNotificationRepo) MarkSeenByIDAndUserID(ctx context.Context, id, userID uuid.UUID) error {
	now := time.Now()
	return r.DB.WithContext(ctx).
		Model(&model.UserNotification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]any{
			"is_seen": true,
			"seen_at": now,
		}).Error
}

func (r *UserNotificationRepo) MarkReadByIDAndUserID(ctx context.Context, id, userID uuid.UUID) error {
	now := time.Now()
	return r.DB.WithContext(ctx).
		Model(&model.UserNotification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]any{
			"is_seen": true,
			"seen_at": now,
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *UserNotificationRepo) MarkAllReadByUserID(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.DB.WithContext(ctx).
		Model(&model.UserNotification{}).
		Where("user_id = ?", userID).
		Where("is_read = ?", false).
		Updates(map[string]any{
			"is_seen": true,
			"seen_at": now,
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *UserNotificationRepo) MarkSSEDeliveredByID(ctx context.Context, id uuid.UUID, channelState int64) error {
	now := time.Now()
	return r.DB.WithContext(ctx).
		Model(&model.UserNotification{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"channel_state": channelState,
			"pushed_at":     now,
		}).Error
}

func (r *UserNotificationRepo) MarkEmailQueuedByIDs(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	now := time.Now()
	return r.DB.WithContext(ctx).
		Model(&model.UserNotification{}).
		Where("id IN ?", ids).
		Updates(map[string]any{
			"emailed_at": now,
		}).Error
}
