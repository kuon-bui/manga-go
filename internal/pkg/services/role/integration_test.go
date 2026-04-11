//go:build integration

package roleservice

import (
	"context"
	"os"
	"testing"
	"time"

	"manga-go/internal/pkg/logger"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	rolerepo "manga-go/internal/pkg/repo/role"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newRoleServiceIntegration(t *testing.T) *RoleService {
	t.Helper()

	dsn := os.Getenv("INTEGRATION_TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("INTEGRATION_TEST_DATABASE_DSN is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})

	ddl := []string{
		`CREATE TABLE roles (
			id uuid PRIMARY KEY,
			name TEXT,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)`,
		`CREATE TABLE permissions (
			id uuid PRIMARY KEY,
			name TEXT,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ
		)`,
		`CREATE TABLE roles_permissions (
			role_id uuid,
			permission_id uuid,
			PRIMARY KEY (role_id, permission_id)
		)`,
	}

	for _, stmt := range ddl {
		if err := tx.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to setup schema: %v", err)
		}
	}

	return &RoleService{
		logger:         logger.NewLogger(),
		roleRepo:       rolerepo.NewRoleRepo(tx),
		permissionRepo: permissionrepo.NewPermissionRepo(tx),
	}
}

func TestRoleServiceIntegrationAssignAndRemovePermission(t *testing.T) {
	s := newRoleServiceIntegration(t)
	ctx := context.Background()
	now := time.Now()

	roleID := uuid.New()
	permOneID := uuid.New()
	permTwoID := uuid.New()

	if err := s.roleRepo.DB.Exec(
		"INSERT INTO roles (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)",
		roleID,
		"editor",
		now,
		now,
	).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	if err := s.roleRepo.DB.Exec(
		"INSERT INTO permissions (id, name, created_at, updated_at) VALUES (?, ?, ?, ?), (?, ?, ?, ?)",
		permOneID,
		"manage_comics",
		now,
		now,
		permTwoID,
		"manage_chapters",
		now,
		now,
	).Error; err != nil {
		t.Fatalf("failed to seed permissions: %v", err)
	}

	assignRes := s.AssignPermissions(ctx, roleID, []uuid.UUID{permOneID, permTwoID})
	if !assignRes.Success {
		t.Fatalf("expected assign permissions success, got: %s", assignRes.Message)
	}

	var assignedCount int64
	if err := s.roleRepo.DB.Table("roles_permissions").Where("role_id = ?", roleID).Count(&assignedCount).Error; err != nil {
		t.Fatalf("failed to count role permissions: %v", err)
	}
	if assignedCount != 2 {
		t.Fatalf("expected 2 assigned permissions, got %d", assignedCount)
	}

	removeRes := s.RemovePermission(ctx, roleID, permOneID)
	if !removeRes.Success {
		t.Fatalf("expected remove permission success, got: %s", removeRes.Message)
	}

	var remainingCount int64
	if err := s.roleRepo.DB.Table("roles_permissions").Where("role_id = ?", roleID).Count(&remainingCount).Error; err != nil {
		t.Fatalf("failed to count remaining permissions: %v", err)
	}
	if remainingCount != 1 {
		t.Fatalf("expected 1 remaining permission, got %d", remainingCount)
	}
}
