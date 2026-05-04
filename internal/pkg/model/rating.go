package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

type Rating struct {
	common.SqlModel
	UserId  uuid.UUID `json:"userId" gorm:"column:user_id"`
	ComicId uuid.UUID `json:"comicId" gorm:"column:comic_id"`
	Score   int       `json:"score" gorm:"column:score"`
	Comment *string   `json:"comment,omitempty" gorm:"column:comment"`

	// Relationships
	User  *User  `json:"user,omitempty" gorm:"foreignKey:UserId"`
	Comic *Comic `json:"comic,omitempty" gorm:"foreignKey:ComicId"`
}

func (Rating) TableName() string {
	return "ratings"
}

func (r *Rating) Fake(f faker.Faker) {
	r.Score = f.IntBetween(1, 5)
	comment := f.Lorem().Sentence(10)
	r.Comment = &comment
}
