package notificationservice

import (
	"context"
	"fmt"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	notificationpkg "manga-go/internal/pkg/notification"
	notificationrequest "manga-go/internal/pkg/request/notification"
	"maps"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (s *NotificationService) CreateTestNotification(ctx context.Context, userID uuid.UUID, req *notificationrequest.CreateTestNotificationRequest) response.Result {
	now := time.Now().UTC()
	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = "Test notification"
	}

	body := strings.TrimSpace(req.Body)
	if body == "" {
		body = fmt.Sprintf("This is a test notification generated at %s", now.Format(time.RFC3339))
	}

	notificationType := req.Type
	if notificationType == "" {
		notificationType = "system.test"
	}

	category := req.Category
	if category == "" {
		category = "system"
	}

	payload := common.JSONMap{}
	maps.Copy(payload, req.Payload)
	if len(payload) == 0 {
		payload["generatedAt"] = now.Format(time.RFC3339)
		payload["source"] = "manual-test"
	}

	tx := s.gormDb.Begin().WithContext(ctx)
	if tx.Error != nil {
		s.logger.Error("Failed to start transaction for test notification", "error", tx.Error)
		return response.ResultErrDb(tx.Error)
	}

	notificationRecord := &model.Notification{
		Type:     notificationType,
		Category: category,
		Title:    title,
		Body:     body,
		Payload:  payload,
	}
	if err := s.notificationRepo.CreateWithTransaction(tx, notificationRecord); err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create test notification", "error", err)
		return response.ResultErrDb(err)
	}

	userNotification := &model.UserNotification{
		NotificationID: notificationRecord.ID,
		UserID:         userID,
		ChannelState:   notificationpkg.ChannelStateSSEQueued,
	}
	if err := s.userNotificationRepo.CreateWithTransaction(tx, userNotification); err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create test user notification", "error", err)
		return response.ResultErrDb(err)
	}

	if err := tx.Commit().Error; err != nil {
		s.logger.Error("Failed to commit test notification transaction", "error", err)
		return response.ResultErrDb(err)
	}

	item, err := s.userNotificationRepo.FindByIDAndUserID(ctx, userNotification.ID, userID)
	if err != nil {
		s.logger.Error("Failed to reload test notification", "error", err)
		return response.ResultErrDb(err)
	}

	message := "Test notification created successfully"
	if err := s.publishRealtime(ctx, []*model.UserNotification{item}, map[uuid.UUID]bool{userID: true}); err != nil {
		s.logger.Error("Failed to publish test notification over SSE", "error", err)
		message = "Test notification created successfully, but SSE delivery failed"
	}

	return response.ResultSuccess(message, s.mapNotificationItem(item))
}
