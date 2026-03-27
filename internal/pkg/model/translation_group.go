package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type TranslationGroup struct {
	common.SqlModel
	Name    string    `json:"name" gorm:"column:name"`
	Slug    string    `json:"slug" gorm:"column:slug"`
	OwnerID uuid.UUID `json:"ownerId" gorm:"column:owner_id"`

	Owner   *User  `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Members []User `json:"members,omitempty" gorm:"foreignKey:TranslationGroupID"`
}

func (TranslationGroup) TableName() string {
	return "translation_groups"
}
