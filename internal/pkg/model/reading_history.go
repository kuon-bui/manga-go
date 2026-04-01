package model

import (
	"manga-go/internal/pkg/common"
	"time"

	"github.com/google/uuid"
)

type ReadingHistory struct {
	common.SqlModel
	UserID     uuid.UUID  `json:"userId" gorm:"column:user_id"`
	ChapterID  uuid.UUID  `json:"chapterId" gorm:"column:chapter_id"`
	ComicID    uuid.UUID  `json:"comicId" gorm:"column:comic_id"`
	LastReadAt *time.Time `json:"lastReadAt" gorm:"column:last_read_at"`

	// Relationships
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Chapter *Chapter `json:"chapter,omitempty" gorm:"foreignKey:ChapterID"`
	Comic   *Comic   `json:"comic,omitempty" gorm:"foreignKey:ComicID"`
}

func (ReadingHistory) TableName() string {
	return "reading_histories"
}
