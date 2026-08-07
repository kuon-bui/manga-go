package roleseeder

import (
	"errors"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/model"
	rolerepo "manga-go/internal/pkg/repo/role"
	seederutil "manga-go/internal/pkg/seeder/util"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// rolePermissions maps role name → the catalog permissions it grants. Names are
// validated against authorization.Catalog when written, so a typo here fails the
// seed rather than producing a role that silently grants nothing.
var rolePermissions = map[string][]string{
	"admin": {
		"comic:read", "comic:write", "comic:delete",
		"chapter:read", "chapter:write", "chapter:delete",
		"user:read", "user:manage",
		"role:manage", "permission:read",
		"audit_log:read",
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
	roleRepo      *rolerepo.RoleRepo
	faker         faker.Faker
	policyManager *authorization.PolicyManager
}

func NewRoleSeeder(
	roleRepo *rolerepo.RoleRepo,
	faker faker.Faker,
	policyManager *authorization.PolicyManager,
) *RoleSeeder {
	return &RoleSeeder{
		roleRepo:      roleRepo,
		faker:         faker,
		policyManager: policyManager,
	}
}

func (s *RoleSeeder) Name() string {
	return "RoleSeeder"
}

// Truncate drops the seeded roles from the database and from the policy engine
// together. Clearing only the table would leave policy rules keyed by role ids
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

	return seederutil.TruncateTables(tx, "roles")
}

func (s *RoleSeeder) Seed(tx *gorm.DB) error {
	// The baseline is what every visitor gets before any role applies: what is
	// publicly readable, and what merely being signed in allows.
	if err := s.policyManager.ReplaceBaselinePolicies(); err != nil {
		return err
	}

	for roleName, permissions := range rolePermissions {
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

		// Re-written on every run so a re-seed repairs a policy engine that has
		// drifted from the roles table.
		if err := s.policyManager.ReplacePermissionsForRole(
			role.ID.String(), permissions, authorization.OrgPlatform,
		); err != nil {
			return err
		}
	}
	return nil
}
