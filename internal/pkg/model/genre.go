package model

import (
	"manga-go/internal/pkg/common"
)

type Genre struct {
	common.SqlModel
	Name string `json:"name" gorm:"column:name"`
}

func (Genre) TableName() string {
	return "genres"
}
