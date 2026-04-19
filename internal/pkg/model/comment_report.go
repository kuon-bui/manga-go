package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type CommentReport struct {
	common.SqlModel
	CommentId uuid.UUID `json:"commentId" gorm:"column:comment_id"`
	UserId    uuid.UUID `json:"userId" gorm:"column:user_id"`
	Reason    string    `json:"reason" gorm:"column:reason"`
	Details   *string   `json:"details,omitempty" gorm:"column:details"`

	Comment *Comment `json:"comment,omitempty" gorm:"foreignKey:CommentId"`
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserId"`
}

func (CommentReport) TableName() string {
	return "comment_reports"
}
