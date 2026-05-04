package pagereactionrepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

func (r *PageReactionRepo) CountByPageIds(ctx context.Context, pageIds []uuid.UUID) (map[uuid.UUID]map[string]int64, error) {
	result := make(map[uuid.UUID]map[string]int64)
	if len(pageIds) == 0 {
		return result, nil
	}

	var counts []struct {
		PageId uuid.UUID
		Type   string
		Count  int64
	}

	if err := r.DB.WithContext(ctx).
		Model(&model.PageReaction{}).
		Where("page_id IN ?", pageIds).
		Where("deleted_at IS NULL").
		Group("page_id, type").
		Select("page_id, type, COUNT(*) as count").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	for _, c := range counts {
		if result[c.PageId] == nil {
			result[c.PageId] = make(map[string]int64)
		}
		result[c.PageId][c.Type] = c.Count
	}

	return result, nil
}

func (r *PageReactionRepo) GetUserReactionsByPageIds(ctx context.Context, pageIds []uuid.UUID, userId uuid.UUID) (map[uuid.UUID]string, error) {
	result := make(map[uuid.UUID]string)
	if len(pageIds) == 0 {
		return result, nil
	}

	var reactions []struct {
		PageId uuid.UUID
		Type   string
	}

	if err := r.DB.WithContext(ctx).
		Model(&model.PageReaction{}).
		Where("page_id IN ?", pageIds).
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL").
		Distinct().
		Select("page_id, type").
		Scan(&reactions).Error; err != nil {
		return nil, err
	}

	for _, reaction := range reactions {
		result[reaction.PageId] = reaction.Type
	}

	return result, nil
}
