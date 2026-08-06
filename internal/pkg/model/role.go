package model

import (
	"manga-go/internal/pkg/common"

	"github.com/jaswdr/faker/v2"
)

// Role carries only the human-facing metadata. Who holds a role, and what it
// grants, live in the policy engine — see authorization.PolicyManager.
type Role struct {
	common.SqlModel
	Name string `json:"name" gorm:"column:name"`

	// Permissions is populated on read from the policy engine; it is never
	// persisted, which is why it carries no gorm mapping.
	Permissions []string `json:"permissions,omitempty" gorm:"-"`
}

func (Role) TableName() string {
	return "roles"
}

func (r *Role) Fake(f faker.Faker) {
	r.Name = common.Slugify(f.Lorem().Sentence(2))
}
