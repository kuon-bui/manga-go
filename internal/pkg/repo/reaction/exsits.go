package reactionrepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

func (r *ReactionRepo) ExistsByCommentIdAndUserId(ctx context.Context, commentId, userId uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&model.Reaction{}).
		Where("comment_id = ? AND user_id = ?", commentId, userId).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
