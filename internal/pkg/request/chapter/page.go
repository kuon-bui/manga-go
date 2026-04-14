package chapterrequest

import "manga-go/internal/pkg/common"

type PageRequest struct {
	PageType common.ContentType `json:"pageType" binding:"omitempty,oneof=image text"`
	ImageURL string             `json:"imageUrl"`
	Content  string             `json:"content"`
}
