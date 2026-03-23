package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type Page struct {
	common.SqlModel
	ChapterID  uuid.UUID `json:"chapterId" gorm:"column:chapter_id"`
	PageNumber int       `json:"pageNumber" gorm:"column:page_number"`
	ImageURL   string    `json:"imageUrl" gorm:"column:image_url"`
}

func (Page) TableName() string {
	return "pages"
}
