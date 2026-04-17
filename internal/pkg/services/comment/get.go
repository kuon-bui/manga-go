package commentservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *CommentService) GetComment(ctx context.Context, id uuid.UUID) response.Result {
	comment, err := s.commentRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, map[string]common.MoreKeyOption{
		"User": {},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Comment")
		}
		s.logger.Error("Failed to get comment", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comment retrieved successfully", comment)
}
