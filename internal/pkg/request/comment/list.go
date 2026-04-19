package commentrequest

import (
	"manga-go/internal/pkg/common"
)

type ListCommentsRequest struct {
	common.Paging

	// Query by chapter (for chapter-level and page-level comments)
	ChapterId string `form:"chapterId"`
	// Query by comic (for comic-level comments on the title detail page)
	ComicId string `form:"comicId"`
	// Filter by page index (within a chapter)
	PageIndex *int `form:"pageIndex"`
}
