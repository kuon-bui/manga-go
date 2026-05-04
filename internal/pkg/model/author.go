package model

import (
	"manga-go/internal/pkg/common"

	"github.com/jaswdr/faker/v2"
)

type Author struct {
	common.SqlModel
	Name string `json:"name" gorm:"column:name"`
}

func (Author) TableName() string {
	return "authors"
}

func (a *Author) Fake(f faker.Faker) {
	a.Name = f.Person().Name()
}
