package userrepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

func (r *UserRepository) AssignRoles(ctx context.Context, userID uuid.UUID, roles []*model.Role) error {
	user := &model.User{}
	user.ID = userID
	return r.DB.WithContext(ctx).Model(user).Association("Roles").Replace(roles)
}

func (r *UserRepository) RemoveRole(ctx context.Context, userID uuid.UUID, role *model.Role) error {
	user := &model.User{}
	user.ID = userID
	return r.DB.WithContext(ctx).Model(user).Association("Roles").Delete(role)
}

func (r *UserRepository) GetRoles(ctx context.Context, userID uuid.UUID) ([]*model.Role, error) {
	user := &model.User{}
	user.ID = userID
	var roles []*model.Role
	err := r.DB.WithContext(ctx).Model(user).Association("Roles").Find(&roles)
	return roles, err
}
