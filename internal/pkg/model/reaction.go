package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

type Reaction struct {
	common.SqlModel
	UserId uuid.UUID `json:"userId" gorm:"column:user_id"`
	Type   string    `json:"type" gorm:"column:type"`

	User *User `json:"user" gorm:"foreignKey:UserId"`
}

func (r *Reaction) Fake(f faker.Faker) {
	types := []string{
		common.ReactionTypeLike,
		common.ReactionTypeLove,
		common.ReactionTypeHaha,
		common.ReactionTypeWow,
		common.ReactionTypeSad,
		common.ReactionTypeAngry,
	}
	r.Type = types[f.IntBetween(0, len(types)-1)]
}
