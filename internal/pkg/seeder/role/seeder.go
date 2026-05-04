package roleseeder

import (
	"errors"
	"manga-go/internal/pkg/model"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	rolerepo "manga-go/internal/pkg/repo/role"
	seederutil "manga-go/internal/pkg/seeder/util"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
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
	faker          faker.Faker
}

func NewRoleSeeder(roleRepo *rolerepo.RoleRepo, permissionRepo *permissionrepo.PermissionRepo, faker faker.Faker) *RoleSeeder {
	return &RoleSeeder{roleRepo: roleRepo, permissionRepo: permissionRepo, faker: faker}
}

func (s *RoleSeeder) Name() string {
	return "RoleSeeder"
}

func (s *RoleSeeder) Truncate(tx *gorm.DB) error {
	return seederutil.TruncateTables(tx, "users_roles", "roles_permissions", "roles")
}

func (s *RoleSeeder) Seed(tx *gorm.DB) error {
	for roleName, permNames := range rolePermissions {
		createdRole := false
		role, err := s.roleRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: roleName}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			role = &model.Role{}
			role.Fake(s.faker)
			role.Name = roleName
			if err := s.roleRepo.CreateWithTransaction(tx, role); err != nil {
				return err
			}
			createdRole = true
		}

		if !createdRole {
			continue
		}

		perms := make([]*model.Permission, 0, len(permNames))
		for _, pn := range permNames {
			perm, err := s.permissionRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: pn}}, nil)
			if err != nil {
				return err
			}
			perms = append(perms, perm)
		}

		if err := s.roleRepo.AssignPermissionsWithTransaction(tx, role.ID, perms); err != nil {
			return err
		}
	}
	return nil
}
