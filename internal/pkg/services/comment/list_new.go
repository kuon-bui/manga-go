package commentservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	commentrequest "manga-go/internal/pkg/request/comment"

	"github.com/google/uuid"
)

func (s *CommentService) ListNewComments(ctx context.Context, req *commentrequest.ListNewCommentsRequest) response.Result {
	comments, total, err := s.commentRepo.FindRecentTopLevelPaginated(ctx, &req.Paging)
	if err != nil {
		s.logger.Error("Failed to list newest comments", "error", err)
		return response.ResultErrDb(err)
	}

	commentIDs := make([]uuid.UUID, 0, len(comments))
	for _, comment := range comments {
		commentIDs = append(commentIDs, comment.ID)
	}

	reactionCounts := make(map[uuid.UUID]map[string]int64)
	if len(commentIDs) > 0 {
		reactionCounts, err = s.reactionRepo.CountByCommentIds(ctx, commentIDs)
		if err != nil {
			s.logger.Error("Failed to fetch reaction counts", "error", err)
			return response.ResultErrDb(err)
		}
	}

	commentResponses := mapNewCommentsToResponses(comments, reactionCounts)
	return response.ResultPaginationData(commentResponses, total, "Newest comments retrieved successfully")
}
