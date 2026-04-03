package commentservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	commentrequest "manga-go/internal/pkg/request/comment"

	"gorm.io/gorm/clause"
)

func (s *CommentService) ListComments(ctx context.Context, req *commentrequest.ListCommentsRequest) response.Result {
	comments, total, err := s.commentRepo.FindPaginated(ctx, []any{
		clause.Eq{Column: "chapter_id", Value: req.ChapterId},
	}, &req.Paging, nil)
	if err != nil {
		s.logger.Error("Failed to list comments", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultPaginationData(comments, total, "Comments retrieved successfully")
}
