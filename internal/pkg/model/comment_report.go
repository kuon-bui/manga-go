package model

import (
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/constant"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

type CommentReport struct {
	common.SqlModel
	CommentId uuid.UUID                    `json:"commentId" gorm:"column:comment_id"`
	UserId    uuid.UUID                    `json:"userId" gorm:"column:user_id"`
	Reason    constant.CommentReportReason `json:"reason" gorm:"column:reason"`
	Details   *string                      `json:"details,omitempty" gorm:"column:details"`

	Comment *Comment `json:"comment,omitempty" gorm:"foreignKey:CommentId"`
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserId"`
}

func (CommentReport) TableName() string {
	return "comment_reports"
}

func (c *CommentReport) Fake(f faker.Faker) {
	reasons := []string{"spam", "abuse", "spoiler", "off-topic"}
	c.Reason = reasons[f.IntBetween(0, len(reasons)-1)]
	details := f.Lorem().Sentence(12)
	c.Details = &details
}
