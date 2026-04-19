package model

import (
	"encoding/json"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/constant"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	UploadedByID       *uuid.UUID              `json:"uploadedById,omitempty" gorm:"column:uploaded_by_id"`
	UploaderId         *uuid.UUID              `json:"uploaderId,omitempty" gorm:"-"`

	// Aggregated stats (computed via subquery, not stored)
	FollowCount  int      `json:"followCount" gorm:"column:follow_count;->"`
	RatingCount  int      `json:"ratingCount" gorm:"column:rating_count;->"`
	ChapterCount int      `json:"chapterCount" gorm:"column:chapter_count;->"`
	AvgRating    *float64 `json:"avgRating" gorm:"column:avg_rating;->"`

	// Relationships
	TranslationGroup *TranslationGroup `json:"translationGroup,omitempty" gorm:"foreignKey:TranslationGroupID"`
	UploadedBy       *User             `json:"uploadedBy,omitempty" gorm:"foreignKey:UploadedByID"`
	Artists          []*Author         `json:"artists" gorm:"many2many:comic_artists;joinForeignKey:ComicID;joinReferences:ArtistID"`
	Authors          []*Author         `json:"authors" gorm:"many2many:comic_authors;"`
	Genres           []*Genre          `json:"genres" gorm:"many2many:comic_genres;"`
	Tags             []*Tag            `json:"tags" gorm:"many2many:comic_tags;"`
	Chapters         []*Chapter        `json:"chapters,omitempty" gorm:"foreignKey:ComicID"`
}

func (Comic) TableName() string {
	return "comics"
}

type comicJSON struct {
	ID                 uuid.UUID               `json:"id"`
	CreatedAt          *time.Time              `json:"createdAt"`
	UpdatedAt          *time.Time              `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt          `json:"deletedAt"`
	Title              string                  `json:"title"`
	Slug               string                  `json:"slug"`
	AlternativeTitles  common.StringSlice      `json:"alternativeTitles"`
	Description        *string                 `json:"description"`
	Thumbnail          *string                 `json:"thumbnail"`
	Banner             *string                 `json:"banner"`
	Type               constant.ComicType      `json:"type"`
	Status             constant.ComicStatus    `json:"status"`
	AgeRating          constant.ComicAgeRating `json:"ageRating"`
	IsPublished        bool                    `json:"isPublished"`
	IsHot              bool                    `json:"isHot"`
	IsFeatured         bool                    `json:"isFeatured"`
	PublishedYear      *int                    `json:"publishedYear"`
	LastChapterAt      *time.Time              `json:"lastChapterAt"`
	TranslationGroupID *uuid.UUID              `json:"translationGroupId,omitempty"`
	UploadedByID       *uuid.UUID              `json:"uploadedById,omitempty"`
	FollowCount        int                     `json:"followCount"`
	RatingCount        int                     `json:"ratingCount"`
	ChapterCount       int                     `json:"chapterCount"`
	AvgRating          *float64                `json:"avgRating"`
	TranslationGroup   *TranslationGroup       `json:"translationGroup,omitempty"`
	UploadedBy         *User                   `json:"uploadedBy,omitempty"`
	Artists            []*Author               `json:"artists"`
	Authors            []*Author               `json:"authors"`
	Genres             []*Genre                `json:"genres"`
	Tags               []*Tag                  `json:"tags"`
	Chapters           []*Chapter              `json:"chapters,omitempty"`
}

func addFileContentPrefix(s *string) *string {
	if s == nil || *s == "" || strings.HasPrefix(*s, "/") || strings.HasPrefix(*s, "http") {
		return s
	}
	v := "/files/content/" + *s
	return &v
}

func (c Comic) MarshalJSON() ([]byte, error) {
	return json.Marshal(comicJSON{
		ID:                 c.ID,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
		DeletedAt:          c.DeletedAt,
		Title:              c.Title,
		Slug:               c.Slug,
		AlternativeTitles:  c.AlternativeTitles,
		Description:        c.Description,
		Thumbnail:          addFileContentPrefix(c.Thumbnail),
		Banner:             addFileContentPrefix(c.Banner),
		Type:               c.Type,
		Status:             c.Status,
		AgeRating:          c.AgeRating,
		IsPublished:        c.IsPublished,
		IsHot:              c.IsHot,
		IsFeatured:         c.IsFeatured,
		PublishedYear:      c.PublishedYear,
		LastChapterAt:      c.LastChapterAt,
		TranslationGroupID: c.TranslationGroupID,
		UploadedByID:       c.UploadedByID,
		FollowCount:        c.FollowCount,
		RatingCount:        c.RatingCount,
		ChapterCount:       c.ChapterCount,
		AvgRating:          c.AvgRating,
		TranslationGroup:   c.TranslationGroup,
		UploadedBy:         c.UploadedBy,
		Artists:            c.Artists,
		Authors:            c.Authors,
		Genres:             c.Genres,
		Tags:               c.Tags,
		Chapters:           c.Chapters,
	})
}
