package comicrequest

import (
	"manga-go/internal/pkg/constant"

	"github.com/google/uuid"
)

type CreateComicRequest struct {
	Title             string                  `json:"title" binding:"required"`
	Slug              string                  `json:"slug" binding:"required"`
	AlternativeTitles []string                `json:"alternativeTitles"`
	Description       *string                 `json:"description"`
	Thumbnail         *string                 `json:"thumbnail"`
	Banner            *string                 `json:"banner"`
	Type              constant.ComicType      `json:"type"`
	AgeRating         constant.ComicAgeRating `json:"ageRating"`
	ArtistId          *uuid.UUID              `json:"artistId"`
	PublishedYear     *int                    `json:"publishedYear"`
	AuthorIDs         []uuid.UUID             `json:"authorIds"`
	GenreSlugs        []string                `json:"genreSlugs"`
	TagSlugs          []string                `json:"tagSlugs"`
}
