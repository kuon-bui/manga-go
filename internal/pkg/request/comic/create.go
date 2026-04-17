package comicrequest

import (
	"manga-go/internal/pkg/constant"
)

type CreateComicRequest struct {
	Title             string                  `json:"title" binding:"required"`
	Slug              string                  `json:"slug" binding:"required"`
	AlternativeTitles []string                `json:"alternativeTitles"`
	Description       *string                 `json:"description"`
	Thumbnail         *string                 `json:"thumbnail"`
	Banner            *string                 `json:"banner"`
	Type              constant.ComicType      `json:"type" binding:"comic_type"`
	AgeRating         constant.ComicAgeRating `json:"ageRating" binding:"required,age_rating"`
	PublishedYear     *int                    `json:"publishedYear"`
	AuthorNames       []string                `json:"authorNames"`
	ArtistNames       []string                `json:"artistNames"`
	GenreSlugs        []string                `json:"genreSlugs"`
	TagSlugs          []string                `json:"tagSlugs"`
}
