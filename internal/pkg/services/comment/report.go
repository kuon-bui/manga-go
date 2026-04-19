package commentservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	commentrequest "manga-go/internal/pkg/request/comment"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *CommentService) ReportComment(ctx context.Context, userId uuid.UUID, commentId uuid.UUID, req *commentrequest.ReportCommentRequest) response.Result {
	// Check if comment exists
	comment, err := s.commentRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: commentId},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Comment")
		}
		s.logger.Error("Failed to find comment", "error", err)
		return response.ResultErrDb(err)
	}

	if comment == nil {
		return response.ResultNotFound("Comment")
	}

	// Check if user has already reported this comment
	_, err = s.commentReportRepo.FindOne(ctx, []any{
		clause.Eq{Column: "comment_id", Value: commentId},
		clause.Eq{Column: "user_id", Value: userId},
		clause.Eq{Column: "deleted_at", Value: nil},
	}, nil)
	if err == nil {
		return response.ResultError("You have already reported this comment")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to check existing report", "error", err)
		return response.ResultErrDb(err)
	}

	// Create report
	report := &model.CommentReport{
		CommentId: commentId,
		UserId:    userId,
		Reason:    req.Reason,
		Details:   req.Details,
	}

	if err := s.commentReportRepo.Create(ctx, report); err != nil {
		s.logger.Error("Failed to create comment report", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comment reported successfully", report)
}
