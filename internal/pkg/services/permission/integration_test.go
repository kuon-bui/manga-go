//go:build integration

package permissionservice

import (
	"context"
	"os"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	permissionrequest "manga-go/internal/pkg/request/permission"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newPermissionServiceIntegration(t *testing.T) (*PermissionService, *gorm.DB) {
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

	if err := tx.Exec(`CREATE TABLE permissions (
		id uuid PRIMARY KEY,
		name TEXT,
		created_at TIMESTAMPTZ,
		updated_at TIMESTAMPTZ,
		deleted_at TIMESTAMPTZ
	)`).Error; err != nil {
		t.Fatalf("failed to setup schema: %v", err)
	}

	s := &PermissionService{
		logger:         logger.NewLogger(),
		permissionRepo: permissionrepo.NewPermissionRepo(tx),
	}

	return s, tx
}

func permissionPaginationTotalFromData(data any) int64 {
	v := reflect.ValueOf(data)
	if !v.IsValid() {
		return -1
	}

	field := v.FieldByName("Total")
	if !field.IsValid() || field.Kind() != reflect.Int64 {
		return -1
	}

	return field.Int()
}

func TestPermissionServiceIntegrationFullFlow(t *testing.T) {
	s, db := newPermissionServiceIntegration(t)
	ctx := context.Background()

	createRes := s.CreatePermission(ctx, &permissionrequest.CreatePermissionRequest{Name: "manage_users"})
	if !createRes.Success {
		t.Fatalf("expected create success, got: %s", createRes.Message)
	}

	listRes := s.ListPermissions(ctx, &common.Paging{Page: 1, Limit: 10})
	if !listRes.Success {
		t.Fatalf("expected list success, got: %s", listRes.Message)
	}
	if total := permissionPaginationTotalFromData(listRes.Data); total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	var permissionID uuid.UUID
	if err := db.Raw("SELECT id FROM permissions WHERE name = ? AND deleted_at IS NULL", "manage_users").Scan(&permissionID).Error; err != nil {
		t.Fatalf("failed to query permission id: %v", err)
	}
	if permissionID == uuid.Nil {
		t.Fatalf("expected persisted permission id")
	}

	updateRes := s.UpdatePermission(ctx, permissionID, &permissionrequest.UpdatePermissionRequest{Name: "manage_system"})
	if !updateRes.Success {
		t.Fatalf("expected update success, got: %s", updateRes.Message)
	}

	deleteRes := s.DeletePermission(ctx, permissionID)
	if !deleteRes.Success {
		t.Fatalf("expected delete success, got: %s", deleteRes.Message)
	}

	notFoundRes := s.UpdatePermission(ctx, permissionID, &permissionrequest.UpdatePermissionRequest{Name: "x"})
	if notFoundRes.Success {
		t.Fatalf("expected not found after soft delete")
	}
	if notFoundRes.Message != "Permission not found" {
		t.Fatalf("unexpected message: %s", notFoundRes.Message)
	}
}
