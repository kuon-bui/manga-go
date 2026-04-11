package comicrequest

import (
	"manga-go/internal/pkg/constant"
)

type UpdateComicRequest struct {
	Title             string                  `json:"title" binding:"required"`
	Slug              string                  `json:"slug" binding:"required"`
	AlternativeTitles []string                `json:"alternativeTitles"`
	Description       *string                 `json:"description"`
	Thumbnail         *string                 `json:"thumbnail"`
	Banner            *string                 `json:"banner"`
	Type              constant.ComicType      `json:"type"`
	AgeRating         constant.ComicAgeRating `json:"ageRating"`
	IsHot             *bool                   `json:"isHot"`
	IsFeatured        *bool                   `json:"isFeatured"`
	PublishedYear     *int                    `json:"publishedYear"`
	AuthorNames       []string                `json:"authorNames"`
	ArtistNames       []string                `json:"artistNames"`
	GenreSlugs        []string                `json:"genreSlugs"`
	TagSlugs          []string                `json:"tagSlugs"`
}
