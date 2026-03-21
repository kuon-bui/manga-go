package model

import (
	"manga-go/internal/pkg/common"
)

type Author struct {
	common.SqlModel
	Name string `json:"name" gorm:"column:name"`
}

func (Author) TableName() string {
	return "authors"
}
