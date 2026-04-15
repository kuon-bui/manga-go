package notificationservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *NotificationService) MarkNotificationSeen(ctx context.Context, userID, notificationID uuid.UUID) response.Result {
	item, err := s.userNotificationRepo.FindByIDAndUserID(ctx, notificationID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Notification")
		}

		s.logger.Error("Failed to find notification for seen update", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.userNotificationRepo.MarkSeenByIDAndUserID(ctx, notificationID, userID); err != nil {
		s.logger.Error("Failed to mark notification as seen", "error", err)
		return response.ResultErrDb(err)
	}

	item.IsSeen = true
	return response.ResultSuccess("Notification marked as seen", s.mapNotificationItem(item))
}

func (s *NotificationService) MarkNotificationRead(ctx context.Context, userID, notificationID uuid.UUID) response.Result {
	item, err := s.userNotificationRepo.FindByIDAndUserID(ctx, notificationID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Notification")
		}

		s.logger.Error("Failed to find notification for read update", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.userNotificationRepo.MarkReadByIDAndUserID(ctx, notificationID, userID); err != nil {
		s.logger.Error("Failed to mark notification as read", "error", err)
		return response.ResultErrDb(err)
	}

	item.IsSeen = true
	item.IsRead = true
	return response.ResultSuccess("Notification marked as read", s.mapNotificationItem(item))
}

func (s *NotificationService) MarkAllNotificationsRead(ctx context.Context, userID uuid.UUID) response.Result {
	if err := s.userNotificationRepo.MarkAllReadByUserID(ctx, userID); err != nil {
		s.logger.Error("Failed to mark all notifications as read", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("All notifications marked as read", nil)
}
