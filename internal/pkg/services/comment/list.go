package commentservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	commentrequest "manga-go/internal/pkg/request/comment"

	"gorm.io/gorm/clause"
)

func (s *CommentService) ListComments(ctx context.Context, req *commentrequest.ListCommentsRequest) response.Result {
	conditions := []any{
		clause.Eq{Column: "chapter_id", Value: req.ChapterId},
		clause.Eq{Column: "parent_id", Value: nil},
	}

	if req.PageIndex != nil {
		conditions = append(conditions, clause.Eq{Column: "page_index", Value: *req.PageIndex})
	}

	comments, total, err := s.commentRepo.FindPaginated(ctx, conditions, &req.Paging, map[string]common.MoreKeyOption{
		"User": {},
	})
	if err != nil {
		s.logger.Error("Failed to list comments", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultPaginationData(comments, total, "Comments retrieved successfully")
}
