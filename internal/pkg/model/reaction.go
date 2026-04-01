package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type Reaction struct {
	common.SqlModel
	UserId    uuid.UUID `json:"userId" gorm:"column:user_id"`
	CommentId uuid.UUID `json:"commentId" gorm:"column:comment_id"`
	Type      string    `json:"type" gorm:"column:type"`

	// Relationships
	User    *User    `json:"user" gorm:"foreignKey:UserId"`
	Comment *Comment `json:"comment" gorm:"foreignKey:CommentId"`
}

func (Reaction) TableName() string {
	return "reactions"
}
