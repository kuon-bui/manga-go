package commentservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	commentrequest "manga-go/internal/pkg/request/comment"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (s *CommentService) HandleReaction(ctx context.Context, user *model.User, commendId uuid.UUID, req *commentrequest.AddReactionRequest) response.Result {
	comment, err := s.commentRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: commendId},
	}, nil)
	if err != nil {
		s.logger.Error("Failed to find comment", "error", err)
		return response.ResultErrDb(err)
	}
	hadReacted, err := s.reactionRepo.ExistsByCommentIdAndUserId(ctx, commendId, user.ID)
	if err != nil {
		s.logger.Error("Failed to check reaction existence", "error", err)
		return response.ResultErrDb(err)
	}

	if hadReacted {
		return s.removeReaction(ctx, user, comment)
	}

	return s.addReaction(ctx, user, comment, req)
}

func (s *CommentService) addReaction(ctx context.Context, user *model.User, comment *model.Comment, req *commentrequest.AddReactionRequest) response.Result {
	reaction := &model.Reaction{
		CommentId: comment.ID,
		UserId:    user.ID,
		Type:      req.Type,
	}

	if err := s.reactionRepo.Create(ctx, reaction); err != nil {
		s.logger.Error("Failed to add reaction", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Reaction added successfully", reaction)
}

func (s *CommentService) removeReaction(ctx context.Context, user *model.User, comment *model.Comment) response.Result {
	if err := s.reactionRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "comment_id", Value: comment.ID},
		clause.Eq{Column: "user_id", Value: user.ID},
	}); err != nil {
		s.logger.Error("Failed to remove reaction", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Reaction removed successfully", nil)
}
