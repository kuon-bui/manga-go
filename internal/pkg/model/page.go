package model

import (
	"encoding/json"
	"manga-go/internal/pkg/common"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Page struct {
	common.SqlModel
	ChapterID  uuid.UUID          `json:"chapterId" gorm:"column:chapter_id"`
	PageNumber int                `json:"pageNumber" gorm:"column:page_number"`
	PageType   common.ContentType `json:"pageType" gorm:"column:page_type"`
	ImageURL   string             `json:"-" gorm:"column:image_url"`
	Content    string             `json:"content" gorm:"column:content"`
}

func (Page) TableName() string {
	return "pages"
}

// pageJSON is used to customize JSON output with full image URL
type pageJSON struct {
	ID         uuid.UUID          `json:"id"`
	CreatedAt  *time.Time         `json:"createdAt"`
	UpdatedAt  *time.Time         `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt     `json:"deletedAt"`
	ChapterID  uuid.UUID          `json:"chapterId"`
	PageNumber int                `json:"pageNumber"`
	PageType   common.ContentType `json:"pageType"`
	ImageURL   string             `json:"imageUrl"`
	Content    string             `json:"content"`
}

func (p Page) MarshalJSON() ([]byte, error) {
	imageURL := p.ImageURL
	if imageURL != "" && !strings.HasPrefix(imageURL, "/") && !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
		imageURL = "/files/content/" + imageURL
	}

	return json.Marshal(pageJSON{
		ID:         p.ID,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		DeletedAt:  p.DeletedAt,
		ChapterID:  p.ChapterID,
		PageNumber: p.PageNumber,
		PageType:   p.PageType,
		ImageURL:   imageURL,
		Content:    p.Content,
	})
}
