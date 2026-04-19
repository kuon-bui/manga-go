package comicrequest

import "manga-go/internal/pkg/common"

type ListComicsRequest struct {
	common.Paging
	TranslationGroupSlug string   `json:"translationGroupSlug" form:"translationGroupSlug"`
	GenreSlugs           []string `json:"genreSlugs" form:"genreSlugs"`
	TagSlugs             []string `json:"tagSlugs" form:"tagSlugs"`
	Search               string   `json:"search" form:"search"`
	SortBy               string   `json:"sortBy" form:"sortBy"`   // lastChapterAt | createdAt | rating | followCount
	Order                string   `json:"order" form:"order"`     // asc | desc (default: desc)
	Status               string   `json:"status" form:"status"`   // ongoing | completed | hiatus | cancelled
}
