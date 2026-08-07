package roleservice

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	authorizationrevision "manga-go/internal/pkg/repo/authorization_revision"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"
	rolerequest "manga-go/internal/pkg/request/role"
	authorizationadmin "manga-go/internal/pkg/services/authorization_admin"
	"manga-go/internal/pkg/testutil"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func newRoleService(t *testing.T, createTables bool) *RoleService {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(testutil.NewSQLiteDSN()), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if createTables {
		testutil.MustSyncSchemas(t, db,
			&testutil.Role{},
			&model.AuthorizationCacheRevision{},
		)
	}
	policyManager := authorization.NewPolicyManager(authorization.PolicyManagerParams{
		Enforcer: testutil.NewInMemoryEnforcer(t),
	})
	roleRepo := rolerepo.NewRoleRepo(db)

	return &RoleService{
		logger:        logger.NewLogger(),
		roleRepo:      roleRepo,
		policyManager: policyManager,
		authAdmin: authorizationadmin.NewService(authorizationadmin.ServiceParams{
			Logger:        logger.NewLogger(),
			RoleRepo:      roleRepo,
			UserRepo:      userrepo.NewUserRepository(db, nil),
			PolicyManager: policyManager,
			Revisions:     authorizationrevision.NewRepo(db),
		}),
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
