package model

import (
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/constant"
	"time"

	"github.com/google/uuid"
)

type Comic struct {
	common.SqlModel
	Title              string                  `json:"title" gorm:"column:title"`
	Slug               string                  `json:"slug" gorm:"column:slug"`
	AlternativeTitles  common.StringSlice      `json:"alternativeTitles" gorm:"column:alternative_titles;type:jsonb"`
	Description        *string                 `json:"description" gorm:"column:description"`
	Thumbnail          *string                 `json:"thumbnail" gorm:"column:thumbnail"`
	Banner             *string                 `json:"banner" gorm:"column:banner"`
	Type               constant.ComicType      `json:"type" gorm:"column:type"`
	Status             constant.ComicStatus    `json:"status" gorm:"column:status"`
	AgeRating          constant.ComicAgeRating `json:"ageRating" gorm:"column:age_rating"`
	IsPublished        bool                    `json:"isPublished" gorm:"column:is_published"`
	IsHot              bool                    `json:"isHot" gorm:"column:is_hot"`
	IsFeatured         bool                    `json:"isFeatured" gorm:"column:is_featured"`
	PublishedYear      *int                    `json:"publishedYear" gorm:"column:published_year"`
	LastChapterAt      *time.Time              `json:"lastChapterAt" gorm:"column:last_chapter_at"`
	TranslationGroupID *uuid.UUID              `json:"translationGroupId,omitempty" gorm:"column:translation_group_id"`

	// Relationships
	TranslationGroup *TranslationGroup `json:"translationGroup,omitempty" gorm:"foreignKey:TranslationGroupID"`
	Artists          []*Author         `json:"artists" gorm:"many2many:comic_artists;joinForeignKey:ComicID;joinReferences:ArtistID"`
	Authors          []*Author         `json:"authors" gorm:"many2many:comic_authors;"`
	Genres           []*Genre          `json:"genres" gorm:"many2many:comic_genres;"`
	Tags             []*Tag            `json:"tags" gorm:"many2many:comic_tags;"`
	Chapters         []*Chapter        `json:"chapters,omitempty" gorm:"foreignKey:ComicID"`
}

func (Comic) TableName() string {
	return "comics"
}
