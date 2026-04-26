package model

import (
	"encoding/json"
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type Page struct {
	common.SqlModel
	ChapterID  uuid.UUID          `json:"chapterId" gorm:"column:chapter_id"`
	PageNumber int                `json:"pageNumber" gorm:"column:page_number"`
	PageType   common.ContentType `json:"pageType" gorm:"column:page_type"`
	ImageURL   string             `json:"imageUrl" gorm:"column:image_url"`
	Content    string             `json:"content" gorm:"column:content"`
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
