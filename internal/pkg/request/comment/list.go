package commentrequest

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type ListCommentsRequest struct {
	common.Paging

	ChapterId uuid.UUID `form:"chapterId" binding:"required"`
	PageIndex int       `form:"pageIndex"`
}
