package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

type Comment struct {
	common.SqlModel
	UserId    uuid.UUID  `json:"userId" gorm:"column:user_id"`
	ChapterId *uuid.UUID `json:"chapterId" gorm:"column:chapter_id"`
	ComicId   uuid.UUID  `json:"comicId" gorm:"column:comic_id"`
	ParentId  *uuid.UUID `json:"parentId,omitempty" gorm:"column:parent_id"`
	PageIndex *int       `json:"pageIndex" gorm:"column:page_index"`
	Content   string     `json:"content" gorm:"column:content"`

	// Relationships
	User    *User      `json:"user" gorm:"foreignKey:UserId"`
	Chapter *Chapter   `json:"chapter" gorm:"foreignKey:ChapterId"`
	Comic   *Comic     `json:"comic" gorm:"foreignKey:ComicId"`
	Parent  *Comment   `json:"parent,omitempty" gorm:"foreignKey:ParentId"`
	Replies []*Comment `json:"replies,omitempty" gorm:"foreignKey:ParentId"`
}

func (Comment) TableName() string {
	return "comments"
}

func (c *Comment) Fake(f faker.Faker) {
	c.Content = f.Lorem().Paragraph(1)
}
