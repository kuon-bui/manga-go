package roleseeder

import (
	"context"
	"manga-go/internal/pkg/model"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	rolerepo "manga-go/internal/pkg/repo/role"

	"gorm.io/gorm/clause"
)

// rolePermissions maps role name → permission names assigned to that role.
var rolePermissions = map[string][]string{
	"admin": {
		"comic:read", "comic:write", "comic:delete",
		"chapter:read", "chapter:write", "chapter:delete",
		"user:read", "user:manage",
		"role:manage",
		"tag:write", "tag:delete",
		"genre:write", "genre:delete",
		"author:write", "author:delete",
		"translation_group:write", "translation_group:delete",
		"comment:delete",
		"rating:delete",
	},
	"translator": {
		"comic:read",
		"chapter:read", "chapter:write",
	},
	"reader": {
		"comic:read",
		"chapter:read",
	},
}

type RoleSeeder struct {
	roleRepo       *rolerepo.RoleRepo
	permissionRepo *permissionrepo.PermissionRepo
}

func NewRoleSeeder(roleRepo *rolerepo.RoleRepo, permissionRepo *permissionrepo.PermissionRepo) *RoleSeeder {
	return &RoleSeeder{roleRepo: roleRepo, permissionRepo: permissionRepo}
}

func (s *RoleSeeder) Name() string {
	return "RoleSeeder"
}

func (s *RoleSeeder) Seed(ctx context.Context) error {
	for roleName, permNames := range rolePermissions {
		role := model.Role{Name: roleName}
		if err := s.roleRepo.DB.WithContext(ctx).
			Where(clause.Eq{Column: "name", Value: roleName}).
			FirstOrCreate(&role).Error; err != nil {
			return err
		}

		perms := make([]*model.Permission, 0, len(permNames))
		for _, pn := range permNames {
			var perm model.Permission
			if err := s.permissionRepo.DB.WithContext(ctx).
				Where(clause.Eq{Column: "name", Value: pn}).
				First(&perm).Error; err != nil {
				return err
			}
			perms = append(perms, &perm)
		}

		if err := s.roleRepo.AssignPermissions(ctx, role.ID, perms); err != nil {
			return err
		}
	}
	return nil
}
