package model

import "github.com/google/uuid"

type PageReaction struct {
	Reaction
	PageId uuid.UUID `json:"pageId" gorm:"column:page_id"`

	Page *Page `json:"page" gorm:"foreignKey:PageId"`
}

func (PageReaction) TableName() string {
	return "page_reactions"
}
