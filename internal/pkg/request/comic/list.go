package comicrequest

import "manga-go/internal/pkg/common"

type ListComicsRequest struct {
	common.Paging
	TranslationGroupSlug string   `json:"translationGroupSlug" form:"translationGroupSlug"`
	GenreSlugs           []string `json:"genreSlugs" form:"genreSlugs"`
	TagSlugs             []string `json:"tagSlugs" form:"tagSlugs"`
	Search               string   `json:"search" form:"search"`
}
