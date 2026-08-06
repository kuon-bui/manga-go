package userservice

import (
	"context"
	"net/http"
	"testing"

	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func newUserServiceWithPolicy(t *testing.T) (*UserService, *authorization.PolicyManager, *gorm.DB) {
	t.Helper()

	db := testutil.NewSQLiteDB(t)
	db.Logger = gormlogger.Discard
	testutil.MustSyncSchemas(t, db, &testutil.User{}, &testutil.Role{})

	pm := authorization.NewPolicyManager(authorization.PolicyManagerParams{
		Enforcer: testutil.NewInMemoryEnforcer(t),
	})

	return &UserService{
		logger:        logger.NewLogger(),
		userRepo:      userrepo.NewUserRepository(db, nil),
		roleRepo:      rolerepo.NewRoleRepo(db),
		policyManager: pm,
	}, pm, db
}

func insertUser(t *testing.T, db *gorm.DB, email string) uuid.UUID {
	t.Helper()

	id := uuid.New()
	if err := db.Exec(
		"INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
		id, "Seed", email, "x",
	).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
	return id
}

func insertRole(t *testing.T, db *gorm.DB, name string) uuid.UUID {
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

func TestAssignRolesWritesToThePolicyEngine(t *testing.T) {
	s, pm, db := newUserServiceWithPolicy(t)
	userID := insertUser(t, db, "a@manga.local")
	roleID := insertRole(t, db, "editor")

	res := s.AssignRoles(context.Background(), userID, []uuid.UUID{roleID})
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected the assignment to succeed, got %d: %s", res.HttpStatus, res.Message)
	}

	roles, err := pm.RolesForUser(userID.String(), authorization.OrgPlatform)
	if err != nil {
		t.Fatalf("failed to read the user's roles: %v", err)
	}
	if len(roles) != 1 || roles[0] != roleID.String() {
		t.Fatalf("expected the user to hold %s, got %v", roleID, roles)
	}
}

func TestAssignRolesReplacesTheExistingSet(t *testing.T) {
	s, pm, db := newUserServiceWithPolicy(t)
	userID := insertUser(t, db, "a@manga.local")
	first := insertRole(t, db, "editor")
	second := insertRole(t, db, "moderator")

	if res := s.AssignRoles(context.Background(), userID, []uuid.UUID{first}); res.HttpStatus != http.StatusOK {
		t.Fatalf("setup failed: %s", res.Message)
	}
	if res := s.AssignRoles(context.Background(), userID, []uuid.UUID{second}); res.HttpStatus != http.StatusOK {
		t.Fatalf("expected the second assignment to succeed, got %s", res.Message)
	}

	roles, err := pm.RolesForUser(userID.String(), authorization.OrgPlatform)
	if err != nil {
		t.Fatalf("failed to read the user's roles: %v", err)
	}
	if len(roles) != 1 || roles[0] != second.String() {
		t.Fatalf("expected only the second role to remain, got %v", roles)
	}
}

func TestRemoveRoleTakesItOutOfThePolicyEngine(t *testing.T) {
	s, pm, db := newUserServiceWithPolicy(t)
	userID := insertUser(t, db, "a@manga.local")
	roleID := insertRole(t, db, "editor")

	if res := s.AssignRoles(context.Background(), userID, []uuid.UUID{roleID}); res.HttpStatus != http.StatusOK {
		t.Fatalf("setup failed: %s", res.Message)
	}

	res := s.RemoveRole(context.Background(), userID, roleID)
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected the removal to succeed, got %d: %s", res.HttpStatus, res.Message)
	}

	roles, err := pm.RolesForUser(userID.String(), authorization.OrgPlatform)
	if err != nil {
		t.Fatalf("failed to read the user's roles: %v", err)
	}
	if len(roles) != 0 {
		t.Fatalf("expected the user to hold no roles, got %v", roles)
	}
}

// The endpoint should answer with role metadata, not bare ids, so the caller
// does not have to resolve every id itself.
func TestGetUserRolesReportsRoleMetadataFromThePolicyEngine(t *testing.T) {
	s, _, db := newUserServiceWithPolicy(t)
	userID := insertUser(t, db, "a@manga.local")
	roleID := insertRole(t, db, "editor")

	if res := s.AssignRoles(context.Background(), userID, []uuid.UUID{roleID}); res.HttpStatus != http.StatusOK {
		t.Fatalf("setup failed: %s", res.Message)
	}

	res := s.GetUserRoles(context.Background(), userID)
	if res.HttpStatus != http.StatusOK {
		t.Fatalf("expected the roles to be returned, got %d", res.HttpStatus)
	}

	roles, ok := res.Data.([]*model.Role)
	if !ok {
		t.Fatalf("expected role metadata in the response, got %T", res.Data)
	}
	if len(roles) != 1 || roles[0].Name != "editor" {
		t.Fatalf("expected the editor role, got %v", roles)
	}
}
