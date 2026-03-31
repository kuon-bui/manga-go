package rolerepo

import (
	"context"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

func (r *RoleRepo) AssignPermissions(ctx context.Context, roleID uuid.UUID, perms []*model.Permission) error {
	role := &model.Role{}
	role.ID = roleID
	return r.DB.WithContext(ctx).Model(role).Association("Permissions").Replace(perms)
}

func (r *RoleRepo) RemovePermission(ctx context.Context, roleID uuid.UUID, perm *model.Permission) error {
	role := &model.Role{}
	role.ID = roleID
	return r.DB.WithContext(ctx).Model(role).Association("Permissions").Delete(perm)
}
