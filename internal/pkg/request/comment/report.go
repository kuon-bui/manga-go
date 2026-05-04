package commentrequest

import "manga-go/internal/pkg/constant"

type ReportCommentRequest struct {
	Reason  constant.CommentReportReason `json:"reason" binding:"required,comment_report_reason"`
	Details *string                      `json:"details" binding:"max=500"`
}
