package pagereactionrepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

func (r *PageReactionRepo) ExistsByPageIdAndUserId(ctx context.Context, pageId, userId uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&model.PageReaction{}).
		Where("page_id = ? AND user_id = ?", pageId, userId).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
