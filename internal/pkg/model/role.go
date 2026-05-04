package model

import (
	"manga-go/internal/pkg/common"

	"github.com/jaswdr/faker/v2"
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

func (r *Role) Fake(f faker.Faker) {
	r.Name = common.Slugify(f.Lorem().Sentence(2))
}
