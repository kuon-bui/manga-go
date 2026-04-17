package commentservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	commentrequest "manga-go/internal/pkg/request/comment"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *CommentService) ListCommentReplies(ctx context.Context, id uuid.UUID, req *commentrequest.ListCommentRepliesRequest) response.Result {
	_, err := s.commentRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Comment")
		}
		s.logger.Error("Failed to find comment", "error", err)
		return response.ResultErrDb(err)
	}

	replies, total, err := s.commentRepo.FindPaginated(ctx, []any{
		clause.Eq{Column: "parent_id", Value: id},
	}, &req.Paging, map[string]common.MoreKeyOption{
		"User": {},
	})
	if err != nil {
		s.logger.Error("Failed to list comment replies", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultPaginationData(replies, total, "Comment replies retrieved successfully")
}
