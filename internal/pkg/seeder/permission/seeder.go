package permissionseeder

import (
	"context"
	"manga-go/internal/pkg/model"
	permissionrepo "manga-go/internal/pkg/repo/permission"

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
	repo *permissionrepo.PermissionRepo
}

func NewPermissionSeeder(repo *permissionrepo.PermissionRepo) *PermissionSeeder {
	return &PermissionSeeder{repo: repo}
}

func (s *PermissionSeeder) Name() string {
	return "PermissionSeeder"
}

func (s *PermissionSeeder) Seed(ctx context.Context) error {
	for _, name := range permissions {
		perm := model.Permission{Name: name}
		if err := s.repo.DB.WithContext(ctx).
			Where(clause.Eq{Column: "name", Value: name}).
			FirstOrCreate(&perm).Error; err != nil {
			return err
		}
	}
	return nil
}
