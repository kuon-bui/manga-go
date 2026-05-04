package permissionseeder

import (
	"errors"
	"manga-go/internal/pkg/model"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	seederutil "manga-go/internal/pkg/seeder/util"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var permissions = []string{
	"comic:read",
	"comic:write",
	"comic:delete",
	"chapter:read",
	"chapter:write",
	"chapter:delete",
	"user:read",
	"user:manage",
	"role:manage",
	"tag:write",
	"tag:delete",
	"genre:write",
	"genre:delete",
	"author:write",
	"author:delete",
	"translation_group:write",
	"translation_group:delete",
	"comment:delete",
	"rating:delete",
}

type PermissionSeeder struct {
	repo  *permissionrepo.PermissionRepo
	faker faker.Faker
}

func NewPermissionSeeder(repo *permissionrepo.PermissionRepo, faker faker.Faker) *PermissionSeeder {
	return &PermissionSeeder{repo: repo, faker: faker}
}

func (s *PermissionSeeder) Name() string {
	return "PermissionSeeder"
}

func (s *PermissionSeeder) Truncate(tx *gorm.DB) error {
	return seederutil.TruncateTables(tx, "roles_permissions", "permissions")
}

func (s *PermissionSeeder) Seed(tx *gorm.DB) error {
	for _, name := range permissions {
		_, err := s.repo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: name}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			perm := &model.Permission{}
			perm.Fake(s.faker)
			perm.Name = name
			if err := s.repo.CreateWithTransaction(tx, perm); err != nil {
				return err
			}
		}
	}
	return nil
}
