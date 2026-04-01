package commentservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	commentrequest "manga-go/internal/pkg/request/comment"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *CommentService) UpdateComment(ctx context.Context, id uuid.UUID, req *commentrequest.UpdateCommentRequest) response.Result {
	comment, err := s.commentRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Comment")
		}
		s.logger.Error("Failed to find comment", "error", err)
		return response.ResultErrDb(err)
	}
	dataUpdate := map[string]any{
		"content": req.Content,
	}
	if req.PageIndex != nil {
		dataUpdate["page_index"] = *req.PageIndex
	}

	if err := s.commentRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, dataUpdate); err != nil {
		s.logger.Error("Failed to update comment", "error", err)
		return response.ResultErrDb(err)
	}

	comment.Content = req.Content
	if req.PageIndex != nil {
		comment.PageIndex = req.PageIndex
	}

	return response.ResultSuccess("Comment updated successfully", comment)
}
