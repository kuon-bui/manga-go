package model

import (
	"manga-go/internal/pkg/common"

	"github.com/jaswdr/faker/v2"
)

type Permission struct {
	common.SqlModel
	Name  string  `json:"name" gorm:"column:name"`
	Roles []*Role `json:"-" gorm:"many2many:roles_permissions;"`
}

func (Permission) TableName() string {
	return "permissions"
}

func (p *Permission) Fake(f faker.Faker) {
	resources := []string{"comic", "chapter", "user", "role", "tag", "genre", "author", "translation_group", "comment", "rating"}
	actions := []string{"read", "write", "delete", "manage"}
	p.Name = resources[f.IntBetween(0, len(resources)-1)] + ":" + actions[f.IntBetween(0, len(actions)-1)]
}
