package model

import (
	"encoding/json"
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

type ReactionCounts struct {
	LIKE  int64 `json:"LIKE"`
	LOVE  int64 `json:"LOVE"`
	HAHA  int64 `json:"HAHA"`
	WOW   int64 `json:"WOW"`
	SAD   int64 `json:"SAD"`
	ANGRY int64 `json:"ANGRY"`
}

type Page struct {
	common.SqlModel
	ChapterID  uuid.UUID          `json:"chapterId" gorm:"column:chapter_id"`
	PageNumber int                `json:"pageNumber" gorm:"column:page_number"`
	PageType   common.ContentType `json:"pageType" gorm:"column:page_type"`
	ImageURL   string             `json:"imageUrl" gorm:"column:image_url"`
	Content    string             `json:"content" gorm:"column:content"`

	ReactionCounts ReactionCounts `json:"reactionCounts" gorm:"-"`
	UserReaction   *string        `json:"userReaction" gorm:"-"`
}

func (Page) TableName() string {
	return "pages"
}

func (p Page) MarshalJSON() ([]byte, error) {
	type alias Page
	temp := alias(p)
	temp.ImageURL = common.AddFileContentPrefix(temp.ImageURL)

	return json.Marshal(temp)
}

func (p *Page) Fake(f faker.Faker) {
	p.PageType = common.ContentTypeImage
	p.ImageURL = f.Internet().URL()
	p.Content = ""

	if f.Bool() {
		p.PageType = common.ContentTypeText
		p.ImageURL = ""
		p.Content = f.Lorem().Paragraph(3)
	}
}
