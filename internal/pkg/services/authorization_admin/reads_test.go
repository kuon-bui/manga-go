package authorizationadmin

import (
	"context"
	"slices"
	"testing"

	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	authorizationrevision "manga-go/internal/pkg/repo/authorization_revision"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type readTestEnv struct {
	service       *Service
	policyManager *authorization.PolicyManager
	db            *gorm.DB
}

func newReadTestEnv(t *testing.T) *readTestEnv {
	t.Helper()
	db := testutil.NewSQLiteDB(t)
	testutil.MustSyncSchemas(
		t,
		db,
		&testutil.Role{},
		&testutil.User{},
		&model.AuthorizationCacheRevision{},
	)
	enforcer := testutil.NewInMemoryEnforcer(t)
	policyManager := authorization.NewPolicyManager(authorization.PolicyManagerParams{Enforcer: enforcer})
	revisions := authorizationrevision.NewRepo(db)

	return &readTestEnv{
		service: &Service{
			logger:        logger.NewLogger(),
			roleRepo:      rolerepo.NewRoleRepo(db),
			userRepo:      userrepo.NewUserRepository(db, nil),
			policyManager: policyManager,
			revisions:     revisions,
		},
		policyManager: policyManager,
		db:            db,
	}
}

func (e *readTestEnv) seedRole(t *testing.T, name string) uuid.UUID {
	t.Helper()
	id := uuid.New()
	if err := e.db.Exec(
		"INSERT INTO roles (id, name, created_at, updated_at) VALUES (?, ?, datetime('now'), datetime('now'))",
		id, name,
	).Error; err != nil {
		t.Fatal(err)
	}
	return id
}

func (e *readTestEnv) seedUser(t *testing.T, name string, email string) uuid.UUID {
	t.Helper()
	id := uuid.New()
	if err := e.db.Exec(
		"INSERT INTO users (id, name, email, created_at, updated_at) VALUES (?, ?, ?, datetime('now'), datetime('now'))",
		id, name, email,
	).Error; err != nil {
		t.Fatal(err)
	}
	return id
}

func TestListUsersSearchesNameAndEmailAndEmbedsRoles(t *testing.T) {
	env := newReadTestEnv(t)
	translatorID := env.seedRole(t, "translator")
	maiID := env.seedUser(t, "Mai Tran", "mai@example.com")
	env.seedUser(t, "Other", "other@example.com")
	if err := env.policyManager.AddRoleForUser(maiID.String(), translatorID.String(), authorization.OrgPlatform); err != nil {
		t.Fatal(err)
	}

	page, err := env.service.ListUsers(context.Background(), ListUsersInput{
		Page: 1, Limit: 20, Search: "MAI@",
	})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Data) != 1 || page.Data[0].Roles[0].Name != "translator" {
		t.Fatalf("unexpected page: %#v", page)
	}
	if page.Data[0].AuthorizationVersion != "g1:u1" {
		t.Fatalf("unexpected authorization version: %q", page.Data[0].AuthorizationVersion)
	}
}

func TestListUsersFiltersByRoleAndKeepsUsersWithoutRolesOut(t *testing.T) {
	env := newReadTestEnv(t)
	translatorID := env.seedRole(t, "translator")
	assignedID := env.seedUser(t, "Assigned", "assigned@example.com")
	env.seedUser(t, "Unassigned", "unassigned@example.com")
	if err := env.policyManager.AddRoleForUser(assignedID.String(), translatorID.String(), authorization.OrgPlatform); err != nil {
		t.Fatal(err)
	}

	page, err := env.service.ListUsers(context.Background(), ListUsersInput{
		Page: 1, Limit: 20, RoleID: &translatorID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Data) != 1 || page.Data[0].ID != assignedID {
		t.Fatalf("unexpected filtered page: %#v", page)
	}
}

func TestListRolesIncludesPermissionsAssignedCountAndVersion(t *testing.T) {
	env := newReadTestEnv(t)
	roleID := env.seedRole(t, "editor")
	for _, userID := range []uuid.UUID{
		env.seedUser(t, "One", "one@example.com"),
		env.seedUser(t, "Two", "two@example.com"),
	} {
		if err := env.policyManager.AddRoleForUser(userID.String(), roleID.String(), authorization.OrgPlatform); err != nil {
			t.Fatal(err)
		}
	}
	if err := env.policyManager.ReplacePermissionsForRole(
		roleID.String(), []string{"comic:write"}, authorization.OrgPlatform,
	); err != nil {
		t.Fatal(err)
	}

	roles, err := env.service.ListRoles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 1 || roles[0].AssignedUserCount != 2 || roles[0].AuthorizationVersion != "g1" {
		t.Fatalf("unexpected roles: %#v", roles)
	}
	if !slices.Equal(roles[0].Permissions, []string{"comic:write"}) {
		t.Fatalf("unexpected permissions: %v", roles[0].Permissions)
	}

	role, err := env.service.GetRole(context.Background(), roleID)
	if err != nil {
		t.Fatal(err)
	}
	if role.ID != roles[0].ID || role.AssignedUserCount != roles[0].AssignedUserCount {
		t.Fatalf("list/detail contract mismatch: list=%#v detail=%#v", roles[0], role)
	}
}
