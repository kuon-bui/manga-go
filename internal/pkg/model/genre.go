package model

import (
	"manga-go/internal/pkg/common"
)

type Genre struct {
	common.SqlModel
	Name        string `json:"name" gorm:"column:name"`
	Slug        string `json:"slug" gorm:"column:slug"`
	Description string `json:"description" gorm:"column:description"`
	Thumbnail   string `json:"thumbnail" gorm:"column:thumbnail"`
}

func (Genre) TableName() string {
	return "genres"
}
