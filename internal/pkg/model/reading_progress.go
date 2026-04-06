package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type ReadingProgress struct {
	common.SqlModel
	UserID        uuid.UUID `json:"userId" gorm:"column:user_id"`
	ComicID       uuid.UUID `json:"comicId" gorm:"column:comic_id"`
	ChapterID     uuid.UUID `json:"chapterId" gorm:"column:chapter_id"`
	ScrollPercent int       `json:"scrollPercent" gorm:"column:scroll_percent"`

	// Relationships
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Comic   *Comic   `json:"comic,omitempty" gorm:"foreignKey:ComicID"`
	Chapter *Chapter `json:"chapter,omitempty" gorm:"foreignKey:ChapterID"`
}

func (ReadingProgress) TableName() string {
	return "reading_progresses"
}
