package commentrequest

import (
	"manga-go/internal/pkg/common"
)

type ListCommentsRequest struct {
	common.Paging

	PageIndex *int   `form:"pageIndex"`
	ChapterId string `form:"chapterId" binding:"required"`
}
