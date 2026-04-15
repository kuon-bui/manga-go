package notificationservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

func (s *NotificationService) ListNotifications(ctx context.Context, userID uuid.UUID, paging *common.Paging, unreadOnly bool, notificationType string) response.Result {
	items, total, err := s.userNotificationRepo.FindPaginatedByUserID(ctx, userID, paging, unreadOnly, notificationType)
	if err != nil {
		s.logger.Error("Failed to list notifications", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultPaginationData(s.mapNotificationItems(items), total, "Notifications retrieved successfully")
}

func (s *NotificationService) mapNotificationItems(items []*model.UserNotification) []NotificationItem {
	res := make([]NotificationItem, 0, len(items))
	for _, item := range items {
		res = append(res, s.mapNotificationItem(item))
	}

	return res
}

func (s *NotificationService) mapNotificationItem(item *model.UserNotification) NotificationItem {
	notification := item.Notification
	if notification == nil {
		notification = &model.Notification{}
	}

	return NotificationItem{
		ID:        item.ID.String(),
		Type:      notification.Type,
		Category:  notification.Category,
		Title:     notification.Title,
		Body:      notification.Body,
		IsSeen:    item.IsSeen,
		IsRead:    item.IsRead,
		CreatedAt: item.CreatedAt,
		Payload:   notification.Payload,
	}
}
