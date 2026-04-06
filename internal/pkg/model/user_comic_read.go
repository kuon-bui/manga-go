package model

import (
	"manga-go/internal/pkg/bitset"
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
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
