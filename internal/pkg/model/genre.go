package model

import (
	"manga-go/internal/pkg/common"

	"github.com/jaswdr/faker/v2"
)

type Genre struct {
	common.SqlModel
	Name        string `json:"name" gorm:"column:name"`
	Slug        string `json:"slug" gorm:"column:slug"`
	Description string `json:"description" gorm:"column:description"`
	Thumbnail   string `json:"thumbnail" gorm:"column:thumbnail"`
}

func (Genre) TableName() string {
	return "genres"
}

func (g *Genre) Fake(f faker.Faker) {
	g.Name = f.Lorem().Word()
	g.Slug = common.Slugify(g.Name)
	g.Description = f.Lorem().Sentence(10)
}
