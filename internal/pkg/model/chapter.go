package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type Chapter struct {
	common.SqlModel
	ComicID     uuid.UUID `json:"comicId" gorm:"column:comic_id"`
	Number      string    `json:"number" gorm:"column:number"`
	ChapterIdx  uint      `json:"-" gorm:"column:chapter_idx;<-:create"`
	Title       string    `json:"title" gorm:"column:title"`
	Slug        string    `json:"slug" gorm:"column:slug"`
	IsPublished bool      `json:"isPublished" gorm:"column:is_published"`

	// Relationships
	Comic *Comic  `json:"comic,omitempty" gorm:"foreignKey:ComicID"`
	Pages []*Page `json:"pages,omitempty" gorm:"foreignKey:ChapterID"`

	IsRead bool `json:"isRead" gorm:"-"`
}

func (Chapter) TableName() string {
	return "chapters"
}
