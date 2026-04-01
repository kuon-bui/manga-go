package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type Comment struct {
	common.SqlModel
	UserId    uuid.UUID `json:"userId" gorm:"column:user_id"`
	ChapterId uuid.UUID `json:"chapterId" gorm:"column:chapter_id"`
	ComicId   uuid.UUID `json:"comicId" gorm:"column:comic_id"`
	PageIndex *int      `json:"pageIndex" gorm:"column:page_index"`
	Content   string    `json:"content" gorm:"column:content"`

	// Relationships
	User    *User    `json:"user" gorm:"foreignKey:UserId"`
	Chapter *Chapter `json:"chapter" gorm:"foreignKey:ChapterId"`
	Comic   *Comic   `json:"comic" gorm:"foreignKey:ComicId"`
}

func (Comment) TableName() string {
	return "comments"
}
