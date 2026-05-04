package model

import (
	"manga-go/internal/pkg/bitset"
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

type UserComicRead struct {
	common.SqlModel
	UserID   uuid.UUID          `json:"userId" gorm:"column:user_id"`
	ComicID  uuid.UUID          `json:"comicId" gorm:"column:comic_id"`
	ReadData *bitset.ReadBitset `json:"readData" gorm:"column:read_data;type:bytea"`

	User  *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Comic *Comic `json:"comic,omitempty" gorm:"foreignKey:ComicID"`
}

func (UserComicRead) TableName() string {
	return "user_comic_reads"
}

func (u *UserComicRead) Fake(f faker.Faker) {
	readData := bitset.NewReadBitset(f.IntBetween(3, 12))
	for index := 0; index < f.IntBetween(1, 4); index++ {
		readData.Mark(index)
	}
	u.ReadData = readData
}
