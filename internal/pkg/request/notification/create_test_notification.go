package notificationrequest

import (
	"manga-go/internal/pkg/common"
	pknotification "manga-go/internal/pkg/notification"
)

type CreateTestNotificationRequest struct {
	Type     pknotification.Type     `json:"type"     binding:"required,notification_type"`
	Category pknotification.Category `json:"category" binding:"required,notification_category"`
	Title    string                  `json:"title"    binding:"required"`
	Body     string                  `json:"body"     binding:"required"`
	Payload  common.JSONMap          `json:"payload"`
}
