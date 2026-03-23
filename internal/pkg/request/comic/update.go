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
	IsActive          *bool                `json:"isActive"`
	IsHot             *bool                `json:"isHot"`
	IsFeatured        *bool                `json:"isFeatured"`
	Artist            *string              `json:"artist"`
	PublishedYear     *int                 `json:"publishedYear"`
	AuthorIDs         []uuid.UUID          `json:"authorIds"`
	GenreIDs          []uuid.UUID          `json:"genreIds"`
	TagIDs            []uuid.UUID          `json:"tagIds"`
}
