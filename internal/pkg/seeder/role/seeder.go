package roleseeder

import (
	"errors"
	"manga-go/internal/pkg/authorization"
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
}

type RoleSeeder struct {
	roleRepo       *rolerepo.RoleRepo
	permissionRepo *permissionrepo.PermissionRepo
	faker          faker.Faker
	policyManager  *authorization.PolicyManager
}

func NewRoleSeeder(
	roleRepo *rolerepo.RoleRepo,
	permissionRepo *permissionrepo.PermissionRepo,
	faker faker.Faker,
	policyManager *authorization.PolicyManager,
) *RoleSeeder {
	return &RoleSeeder{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		faker:          faker,
		policyManager:  policyManager,
	}
}

func (s *RoleSeeder) Name() string {
	return "RoleSeeder"
}

// Truncate drops the seeded roles from the database and from the policy engine
// together. Clearing only the tables would leave policy rules keyed by role ids
// that no longer exist.
func (s *RoleSeeder) Truncate(tx *gorm.DB) error {
	for roleName := range rolePermissions {
		role, err := s.roleRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: roleName}}, nil)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return err
		}
		if err := s.policyManager.RemoveRole(role.ID.String(), authorization.OrgPlatform); err != nil {
			return err
		}
	}

	return seederutil.TruncateTables(tx, "users_roles", "roles_permissions", "roles")
}

func (s *RoleSeeder) Seed(tx *gorm.DB) error {
	for roleName, permNames := range rolePermissions {
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
		}

		perms := make([]*model.Permission, 0, len(permNames))
		rules := make([]authorization.PermissionRule, 0, len(permNames))
		for _, pn := range permNames {
			perm, err := s.permissionRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: pn}}, nil)
			if err != nil {
				return err
			}
			perms = append(perms, perm)
			rules = append(rules, authorization.PermissionRule{ID: perm.ID.String(), Name: perm.Name})
		}

		if err := s.roleRepo.AssignPermissionsWithTransaction(tx, role.ID, perms); err != nil {
			return err
		}

		// A role means nothing until the policy engine knows about it: without
		// this the seeded admin exists in the database but is denied everything.
		if err := s.policyManager.ReplacePermissionsForRole(role.ID.String(), rules, authorization.OrgPlatform); err != nil {
			return err
		}
	}
	return nil
}
