package userrequest

type UpdateUserConfigRequest struct {
	SeenNotificationCenter             *bool `json:"seenNotificationCenter"`
	EnableSSENotifications             *bool `json:"enableSseNotifications"`
	EnableEmailNotifications           *bool `json:"enableEmailNotifications"`
	EnableComicNewChapterNotifications *bool `json:"enableComicNewChapterNotifications"`
	EnableCommentReplyNotifications    *bool `json:"enableCommentReplyNotifications"`
	EnableMentionNotifications         *bool `json:"enableMentionNotifications"`
	EnableSystemAnnouncements          *bool `json:"enableSystemAnnouncements"`
}
