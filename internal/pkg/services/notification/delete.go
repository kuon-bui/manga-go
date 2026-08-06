package notificationservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *NotificationService) DeleteNotification(ctx context.Context, userID, notificationID uuid.UUID) response.Result {
	_, err := s.userNotificationRepo.FindByIDAndUserID(ctx, notificationID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Notification")
		}
		s.logger.Error("Failed to find notification for deletion", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.userNotificationRepo.DeleteByIDAndUserID(ctx, notificationID, userID); err != nil {
		s.logger.Error("Failed to delete notification", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Notification deleted successfully", nil)
}

func (s *NotificationService) DeleteAllNotifications(ctx context.Context, userID uuid.UUID) response.Result {
	if err := s.userNotificationRepo.DeleteAllByUserID(ctx, userID); err != nil {
		s.logger.Error("Failed to delete all notifications", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("All notifications deleted successfully", nil)
}
