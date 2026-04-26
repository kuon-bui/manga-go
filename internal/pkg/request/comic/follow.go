package comicrequest

import "manga-go/internal/pkg/constant"

type FollowComicRequest struct {
	FollowStatus constant.FollowStatus `json:"followStatus" binding:"required,follow_status"`
}
