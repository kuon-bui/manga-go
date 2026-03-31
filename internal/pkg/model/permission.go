package model

import (
	"manga-go/internal/pkg/common"
)

type Permission struct {
	common.SqlModel
	Name  string  `json:"name" gorm:"column:name"`
	Roles []*Role `json:"-" gorm:"many2many:roles_permissions;"`
}

func (Permission) TableName() string {
	return "permissions"
}
