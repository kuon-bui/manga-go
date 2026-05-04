package roleservice

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	rolerepo "manga-go/internal/pkg/repo/role"
	rolerequest "manga-go/internal/pkg/request/role"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newRoleService(t *testing.T, createTables bool) *RoleService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if createTables {
		err = db.Exec(`
			CREATE TABLE roles (
				id TEXT PRIMARY KEY,
				name TEXT,
				created_at DATETIME,
				updated_at DATETIME,
				deleted_at DATETIME
			)
		`).Error
		if err != nil {
			t.Fatalf("failed to create roles table: %v", err)
		}

		err = db.Exec(`
			CREATE TABLE permissions (
				id TEXT PRIMARY KEY,
				name TEXT,
				created_at DATETIME,
				updated_at DATETIME,
				deleted_at DATETIME
			)
		`).Error
		if err != nil {
			t.Fatalf("failed to create permissions table: %v", err)
		}

		err = db.Exec(`
			CREATE TABLE roles_permissions (
				role_id TEXT,
				permission_id TEXT
			)
		`).Error
		if err != nil {
			t.Fatalf("failed to create roles_permissions table: %v", err)
		}
	}

	return &RoleService{
		logger:         logger.NewLogger(),
		roleRepo:       rolerepo.NewRoleRepo(db),
		permissionRepo: permissionrepo.NewPermissionRepo(db),
	}
}

func rolePaginationTotal(data any) int64 {
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

func TestListRolesReturnsEmptyPagination(t *testing.T) {
	t.Parallel()

	s := newRoleService(t, true)
	res := s.ListRoles(context.Background(), &common.Paging{Page: 1, Limit: 10})

	if !res.Success {
		t.Fatalf("expected success result")
	}
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.HttpStatus)
	}
	if res.Message != "Roles retrieved successfully" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
	if total := rolePaginationTotal(res.Data); total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
}

func TestGetRoleReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newRoleService(t, true)
	res := s.GetRole(context.Background(), uuid.New())

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Role not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestUpdateRoleReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newRoleService(t, true)
	res := s.UpdateRole(context.Background(), uuid.New(), &rolerequest.UpdateRoleRequest{Name: "editor"})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Role not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestDeleteRoleReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newRoleService(t, true)
	res := s.DeleteRole(context.Background(), uuid.New())

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Role not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestCreateRoleReturnsDbErrorWhenTableMissing(t *testing.T) {
	t.Parallel()

	s := newRoleService(t, false)
	res := s.CreateRole(context.Background(), &rolerequest.CreateRoleRequest{Name: "admin"})

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

func TestAssignPermissionsReturnsRoleNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newRoleService(t, true)
	res := s.AssignPermissions(context.Background(), uuid.New(), []uuid.UUID{uuid.New()})

	if res.Success {
		t.Fatalf("expected failure result")
	}
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.HttpStatus)
	}
	if res.Message != "Role not found" {
		t.Fatalf("unexpected message: %s", res.Message)
	}
}

func TestAssignPermissionsReturnsPermissionNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newRoleService(t, true)
	roleID := uuid.New()
	if err := s.roleRepo.DB.Exec("INSERT INTO roles (id, name) VALUES (?, ?)", roleID.String(), "editor").Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	res := s.AssignPermissions(context.Background(), roleID, []uuid.UUID{uuid.New()})

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

func TestRemovePermissionReturnsPermissionNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	s := newRoleService(t, true)
	roleID := uuid.New()
	if err := s.roleRepo.DB.Exec("INSERT INTO roles (id, name) VALUES (?, ?)", roleID.String(), "editor").Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}

	res := s.RemovePermission(context.Background(), roleID, uuid.New())

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
