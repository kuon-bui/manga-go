package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type Reaction struct {
	common.SqlModel
	UserId uuid.UUID `json:"userId" gorm:"column:user_id"`
	Type   string    `json:"type" gorm:"column:type"`

	User *User `json:"user" gorm:"foreignKey:UserId"`
}
