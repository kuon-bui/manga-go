package userrepo

import (
	"context"
	"strings"

	"manga-go/internal/pkg/model"
	userrequest "manga-go/internal/pkg/request/user"

	"github.com/google/uuid"
)

func (r *UserRepository) ListAuthorizationUsers(
	ctx context.Context,
	input *userrequest.ListAuthorizationUsersRequest,
	roleUserIDs []uuid.UUID,
) ([]*model.User, int64, error) {
	input.Fulfill()
	if input.RoleID != nil && len(roleUserIDs) == 0 {
		return []*model.User{}, 0, nil
	}

	query := r.DB.WithContext(ctx).Model(&model.User{})
	if search := strings.TrimSpace(strings.ToLower(input.Search)); search != "" {
		pattern := "%" + search + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", pattern, pattern)
	}
	if input.RoleID != nil {
		query = query.Where("id IN ?", roleUserIDs)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []*model.User
	if err := query.
		Order("created_at DESC").
		Limit(input.GetLimit()).
		Offset(input.GetOffset()).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
