package permissionservice

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	permissionrequest "manga-go/internal/pkg/request/permission"
	"manga-go/internal/pkg/testutil"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newPermissionService(t *testing.T, createTable bool) *PermissionService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(testutil.NewSQLiteDSN()), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if createTable {
		testutil.MustSyncSchemas(t, db, &testutil.Permission{})
	}

	return &PermissionService{
		logger:         logger.NewLogger(),
		permissionRepo: permissionrepo.NewPermissionRepo(db),
	}
}

func permissionPaginationTotal(data any) int64 {
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

func TestListPermissionsReturnsEmptyPagination(t *testing.T) {
	t.Parallel()

	s := newPermissionService(t, true)
	res := s.ListPermissions(context.Background(), &common.Paging{Page: 1, Limit: 10})

	if !res.Success {
		t.Fatalf("expected success result")
	}
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.HttpStatus)
	}
	if res.Message != "Permissions retrieved successfully" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if total := permissionPaginationTotal(res.Data); total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
}

func TestUpdatePermissionReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newPermissionService(t, true)
	res := s.UpdatePermission(context.Background(), uuid.New(), &permissionrequest.UpdatePermissionRequest{Name: "admin"})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Permission not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestDeletePermissionReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newPermissionService(t, true)
	res := s.DeletePermission(context.Background(), uuid.New())

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Permission not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestCreatePermissionReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newPermissionService(t, false)
	res := s.CreatePermission(context.Background(), &permissionrequest.CreatePermissionRequest{Name: "manage_users"})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", res.HttpStatus)
	}
	if res.Message != "database error" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if res.Error == nil {
		t.Fatalf("expected non-nil error")
	}
}
