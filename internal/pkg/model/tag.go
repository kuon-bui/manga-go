package model

import (
	"manga-go/internal/pkg/common"

	"github.com/jaswdr/faker/v2"
)

type Tag struct {
	common.SqlModel
	Name string `json:"name" gorm:"column:name"`
	Slug string `json:"slug" gorm:"column:slug"`
}

func (Tag) TableName() string {
	return "tags"
}

func (t *Tag) Fake(f faker.Faker) {
	t.Name = f.Lorem().Word()
	t.Slug = common.Slugify(t.Name)
}
