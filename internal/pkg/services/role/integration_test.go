//go:build integration

package roleservice

import (
	"context"
	"testing"
	"time"

	"manga-go/internal/pkg/logger"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	rolerepo "manga-go/internal/pkg/repo/role"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
)

func newRoleServiceIntegration(t *testing.T) *RoleService {
	t.Helper()

	db := testutil.NewSQLiteDB(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})

	testutil.MustSyncSchemas(t, tx,
		&testutil.Role{},
		&testutil.Permission{},
		&testutil.RolePermission{},
	)

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
