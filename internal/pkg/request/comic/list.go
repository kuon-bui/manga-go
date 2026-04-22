package comicrequest

import (
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/constant"
)

type ListComicsRequest struct {
	common.Paging
	TranslationGroupID string               `json:"translationGroupSlug" form:"translationGroupSlug"`
	GenreSlugs         []string             `json:"genreSlugs" form:"genreSlugs"`
	TagSlugs           []string             `json:"tagSlugs" form:"tagSlugs"`
	Search             string               `json:"search" form:"search"`
	SortBy             string               `json:"sortBy" form:"sortBy" binding:"comic_sort_by"` // lastChapterAt | createdAt | rating | followCount
	Order              string               `json:"order" form:"order" binding:"order_check"`     // asc | desc (default: desc)
	Status             constant.ComicStatus `json:"status" form:"status" binding:"comic_status"`  // ongoing | completed | hiatus | cancelled
}
