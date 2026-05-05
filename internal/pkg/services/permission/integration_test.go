//go:build integration

package permissionservice

import (
	"context"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	permissionrequest "manga-go/internal/pkg/request/permission"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newPermissionServiceIntegration(t *testing.T) (*PermissionService, *gorm.DB) {
	t.Helper()

	db := testutil.NewSQLiteDB(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})

	testutil.MustSyncSchemas(t, tx, &testutil.Permission{})

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

	permissionID := testutil.MustReadUUID(t, db, "SELECT id FROM permissions WHERE name = ? AND deleted_at IS NULL", "manage_users")
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
