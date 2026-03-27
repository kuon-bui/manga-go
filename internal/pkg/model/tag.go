package model

import (
	"manga-go/internal/pkg/common"
)

type Tag struct {
	common.SqlModel
	Name string `json:"name" gorm:"column:name"`
	Slug string `json:"slug" gorm:"column:slug"`
}

func (Tag) TableName() string {
	return "tags"
}
