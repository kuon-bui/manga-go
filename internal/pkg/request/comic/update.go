package comicrequest

import (
	"manga-go/internal/pkg/constant"

	"github.com/google/uuid"
)

type UpdateComicRequest struct {
	Title             string               `json:"title" binding:"required"`
	Slug              string               `json:"slug" binding:"required"`
	AlternativeTitles []string             `json:"alternativeTitles"`
	Description       *string              `json:"description"`
	Thumbnail         *string              `json:"thumbnail"`
	Banner            *string              `json:"banner"`
	Type              constant.ComicType   `json:"type"`
	Status            constant.ComicStatus `json:"status"`
	IsHot             *bool                `json:"isHot"`
	IsFeatured        *bool                `json:"isFeatured"`
	PublishedYear     *int                 `json:"publishedYear"`
	AuthorIDs         []uuid.UUID          `json:"authorIds"`
	GenreSlugs        []string             `json:"genreSlugs"`
	TagSlugs          []string             `json:"tagSlugs"`
}
