package notificationrequest

import (
	"manga-go/internal/pkg/common"
	pknotification "manga-go/internal/pkg/notification"
)

type ListNotificationsRequest struct {
	common.Paging
	UnreadOnly bool                `form:"unreadOnly"`
	Type       pknotification.Type `form:"type" binding:"notification_type"`
}
