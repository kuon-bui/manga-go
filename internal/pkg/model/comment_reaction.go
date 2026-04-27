package model

import "github.com/google/uuid"

type CommentReaction struct {
	Reaction
	CommentId uuid.UUID `json:"commentId" gorm:"column:comment_id"`

	Comment *Comment `json:"comment" gorm:"foreignKey:CommentId"`
}

func (CommentReaction) TableName() string {
	return "comment_reactions"
}
