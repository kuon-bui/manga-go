package model

import (
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/constant"

	"github.com/google/uuid"
)

type ComicFollow struct {
	common.SqlModel
	UserID       uuid.UUID             `json:"userId" gorm:"column:user_id"`
	ComicID      uuid.UUID             `json:"comicId" gorm:"column:comic_id"`
	FollowStatus constant.FollowStatus `json:"followStatus" gorm:"column:follow_status"`

	User  *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Comic *Comic `json:"comic,omitempty" gorm:"foreignKey:ComicID"`
}

func (ComicFollow) TableName() string {
	return "comic_follows"
}
