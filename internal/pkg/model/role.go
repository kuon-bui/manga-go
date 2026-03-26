package model

import (
	"manga-go/internal/pkg/common"
)

type Role struct {
	common.SqlModel
	Name        string        `json:"name" gorm:"column:name"`
	Permissions []*Permission `json:"permissions,omitempty" gorm:"many2many:roles_permissions;"`
	Users       []*User       `json:"-" gorm:"many2many:users_roles;"`
}

func (Role) TableName() string {
	return "roles"
}
