package roleservice

import (
	"context"
	"net/http"
	"testing"

	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	authorizationrevision "manga-go/internal/pkg/repo/authorization_revision"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"
	rolerequest "manga-go/internal/pkg/request/role"
	authorizationadmin "manga-go/internal/pkg/services/authorization_admin"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func newRoleServiceWithPolicy(t *testing.T) (*RoleService, *authorization.PolicyManager, *gorm.DB) {
	t.Helper()

	db := testutil.NewSQLiteDB(t)
	db.Logger = gormlogger.Discard
	testutil.MustSyncSchemas(t, db, &testutil.Role{}, &model.AuthorizationCacheRevision{})

	pm := authorization.NewPolicyManager(authorization.PolicyManagerParams{
		Enforcer: testutil.NewInMemoryEnforcer(t),
	})

	roleRepo := rolerepo.NewRoleRepo(db)
	authAdmin := authorizationadmin.NewService(authorizationadmin.ServiceParams{
		Logger:        logger.NewLogger(),
		RoleRepo:      roleRepo,
		UserRepo:      userrepo.NewUserRepository(db, nil),
		PolicyManager: pm,
		Revisions:     authorizationrevision.NewRepo(db),
	})

	return &RoleService{
		logger:        logger.NewLogger(),
		roleRepo:      roleRepo,
		policyManager: pm,
		authAdmin:     authAdmin,
	}, pm, db
}

func seedRole(t *testing.T, db *gorm.DB, name string) uuid.UUID {
	t.Helper()

	id := uuid.New()
	if err := db.Exec(
		"INSERT INTO roles (id, name, created_at, updated_at) VALUES (?, ?, datetime('now'), datetime('now'))",
		id, name,
	).Error; err != nil {
		t.Fatalf("failed to seed role: %v", err)
	}
	return id
}

func TestAssignPermissionsStoresGrantsInThePolicyEngine(t *testing.T) {
	s, pm, db := newRoleServiceWithPolicy(t)
	roleID := seedRole(t, db, "editor")

	res := s.AssignPermissions(context.Background(), roleID, &rolerequest.AssignPermissionRequest{
		Permissions: []string{"comic:read", "comic:write"},
	})
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected the grant to succeed, got %d: %s", res.HttpStatus, res.Message)
	}

	names, err := pm.PermissionNamesForRole(roleID.String(), authorization.OrgPlatform)
	if err != nil {
		t.Fatalf("failed to read back the grants: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected comic:read and comic:write, got %v", names)
	}
}

func TestAssignPermissionsRejectsUnknownName(t *testing.T) {
	s, pm, db := newRoleServiceWithPolicy(t)
	roleID := seedRole(t, db, "editor")

	res := s.AssignPermissions(context.Background(), roleID, &rolerequest.AssignPermissionRequest{
		Permissions: []string{"comic:read", "comic:frobnicate"},
	})
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected 400 for an unknown permission, got %d", res.HttpStatus)
	}

	names, err := pm.PermissionNamesForRole(roleID.String(), authorization.OrgPlatform)
	if err != nil {
		t.Fatalf("failed to read back the grants: %v", err)
	}
	if len(names) != 0 {
		t.Fatalf("a rejected request must not have written anything, got %v", names)
	}
}

func TestAssignPermissionsReturnsNotFoundForMissingRole(t *testing.T) {
	s, _, _ := newRoleServiceWithPolicy(t)

	res := s.AssignPermissions(context.Background(), uuid.New(), &rolerequest.AssignPermissionRequest{
		Permissions: []string{"comic:read"},
	})
	if res.Message != "Role not found" {
		t.Fatalf("expected a not-found result, got %q", res.Message)
	}
}

func TestRemovePermissionRevokesOnlyThatGrant(t *testing.T) {
	s, pm, db := newRoleServiceWithPolicy(t)
	roleID := seedRole(t, db, "editor")

	if res := s.AssignPermissions(context.Background(), roleID, &rolerequest.AssignPermissionRequest{
		Permissions: []string{"comic:read", "role:manage"},
	}); res.HttpStatus != http.StatusOK {
		t.Fatalf("setup failed: %s", res.Message)
	}

	res := s.RemovePermission(context.Background(), roleID, "comic:read")
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected the revoke to succeed, got %d: %s", res.HttpStatus, res.Message)
	}

	names, err := pm.PermissionNamesForRole(roleID.String(), authorization.OrgPlatform)
	if err != nil {
		t.Fatalf("failed to read back the grants: %v", err)
	}
	if len(names) != 1 || names[0] != "role:manage" {
		t.Fatalf("expected only role:manage to remain, got %v", names)
	}
}

func TestRemovePermissionRejectsUnknownName(t *testing.T) {
	s, _, db := newRoleServiceWithPolicy(t)
	roleID := seedRole(t, db, "editor")

	res := s.RemovePermission(context.Background(), roleID, "comic:frobnicate")
	if res.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected 400 for an unknown permission, got %d", res.HttpStatus)
	}
}

// The role endpoint has to report grants, and the policy engine is now the only
// place they exist.
func TestGetRoleReportsPermissionsFromThePolicyEngine(t *testing.T) {
	s, _, db := newRoleServiceWithPolicy(t)
	roleID := seedRole(t, db, "editor")

	if res := s.AssignPermissions(context.Background(), roleID, &rolerequest.AssignPermissionRequest{
		Permissions: []string{"comic:write"},
	}); res.HttpStatus != http.StatusOK {
		t.Fatalf("setup failed: %s", res.Message)
	}

	res := s.GetRole(context.Background(), roleID)
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected the role to be returned, got %d", res.HttpStatus)
	}

	role, ok := res.Data.(*authorizationadmin.RoleAccessSummary)
	if !ok {
		t.Fatalf("expected a role in the response, got %T", res.Data)
	}
	if len(role.Permissions) != 1 || role.Permissions[0] != "comic:write" {
		t.Fatalf("expected the role to report comic:write, got %v", role.Permissions)
	}
}
