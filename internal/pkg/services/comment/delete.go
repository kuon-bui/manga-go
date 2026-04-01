package commentservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *CommentService) DeleteComment(ctx context.Context, id uuid.UUID) response.Result {
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

	if err := s.commentRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}); err != nil {
		s.logger.Error("Failed to delete comment", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comment deleted successfully", comment)
}
