package reactionrepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

type ReactionCountByType struct {
	Type  string
	Count int64
}

// CountByCommentIds returns reaction counts grouped by comment_id and type
// Returns map[commentId]map[reactionType]count
func (r *ReactionRepo) CountByCommentIds(ctx context.Context, commentIds []uuid.UUID) (map[uuid.UUID]map[string]int64, error) {
	var counts []struct {
		CommentId uuid.UUID
		Type      string
		Count     int64
	}

	if err := r.DB.WithContext(ctx).
		Model(&model.Reaction{}).
		Where("comment_id IN ?", commentIds).
		Where("deleted_at IS NULL").
		Group("comment_id, type").
		Select("comment_id, type, COUNT(*) as count").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]map[string]int64)
	for _, c := range counts {
		if result[c.CommentId] == nil {
			result[c.CommentId] = make(map[string]int64)
		}
		result[c.CommentId][c.Type] = c.Count
	}

	return result, nil
}

// GetUserReactionsByCommentIds returns user's reaction for each comment
// Returns map[commentId]reactionType (empty string if no reaction)
func (r *ReactionRepo) GetUserReactionsByCommentIds(ctx context.Context, commentIds []uuid.UUID, userId uuid.UUID) (map[uuid.UUID]string, error) {
	var reactions []struct {
		CommentId uuid.UUID
		Type      string
	}

	if err := r.DB.WithContext(ctx).
		Model(&model.Reaction{}).
		Where("comment_id IN ?", commentIds).
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL").
		Distinct().
		Select("comment_id, type").
		Scan(&reactions).Error; err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]string)
	for _, reaction := range reactions {
		result[reaction.CommentId] = reaction.Type
	}

	return result, nil
}
