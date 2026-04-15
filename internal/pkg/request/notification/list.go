package notificationrequest

import "manga-go/internal/pkg/common"

type ListNotificationsRequest struct {
	common.Paging
	UnreadOnly bool   `form:"unreadOnly"`
	Type       string `form:"type"`
}
