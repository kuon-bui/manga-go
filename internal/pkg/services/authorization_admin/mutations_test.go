package authorizationadmin

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	authorizationaudit "manga-go/internal/pkg/repo/authorization_audit"
	authorizationrevision "manga-go/internal/pkg/repo/authorization_revision"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type failingAuditStore struct {
	base *authorizationaudit.Repo
	err  error
}

func (f *failingAuditStore) AppendTx(*gorm.DB, *model.AuthorizationAuditLog) error {
	return f.err
}

func (f *failingAuditStore) List(
	ctx context.Context,
	input authorizationaudit.ListInput,
) ([]*model.AuthorizationAuditLog, int64, error) {
	return f.base.List(ctx, input)
}

type mutationTestEnv struct {
	service       *Service
	policyManager *authorization.PolicyManager
	revisions     *authorizationrevision.Repo
	auditRepo     *authorizationaudit.Repo
	cache         *fakeProfileCache
	db            *gorm.DB
}

func newMutationTestEnv(t *testing.T, auditStore AuditStore) *mutationTestEnv {
	t.Helper()
	db := testutil.NewSQLiteDB(t)
	testutil.MustSyncSchemas(
		t,
		db,
		&testutil.User{},
		&testutil.Role{},
		&model.AuthorizationCacheRevision{},
		&testutil.AuthorizationAuditLog{},
	)
	enforcer := testutil.NewInMemoryEnforcer(t)
	policyManager := authorization.NewPolicyManager(authorization.PolicyManagerParams{Enforcer: enforcer})
	revisions := authorizationrevision.NewRepo(db)
	auditRepo := authorizationaudit.NewRepo(db)
	if auditStore == nil {
		auditStore = auditRepo
	}
	cache := &fakeProfileCache{}
	service := &Service{
		logger:        logger.NewLogger(),
		db:            db,
		roleRepo:      rolerepo.NewRoleRepo(db),
		userRepo:      userrepo.NewUserRepository(db, nil),
		policyManager: policyManager,
		authorizer:    authorization.NewAuthorizer(enforcer),
		revisions:     revisions,
		auditRepo:     auditStore,
		cache:         cache,
		locker:        NewMutexMutationLocker(),
	}
	return &mutationTestEnv{
		service:       service,
		policyManager: policyManager,
		revisions:     revisions,
		auditRepo:     auditRepo,
		cache:         cache,
		db:            db,
	}
}

func (e *mutationTestEnv) seedUser(t *testing.T, name string) *model.User {
	t.Helper()
	user := &model.User{
		SqlModel: common.SqlModel{ID: uuid.New()},
		Name:     name,
		Email:    name + "@example.com",
	}
	if err := e.db.Create(user).Error; err != nil {
		t.Fatal(err)
	}
	return user
}

func (e *mutationTestEnv) seedRole(t *testing.T, name string, permissions ...string) *model.Role {
	t.Helper()
	role := &model.Role{SqlModel: common.SqlModel{ID: uuid.New()}, Name: name}
	if err := e.db.Create(role).Error; err != nil {
		t.Fatal(err)
	}
	if len(permissions) > 0 {
		if err := e.policyManager.ReplacePermissionsForRole(role.ID.String(), permissions, authorization.OrgPlatform); err != nil {
			t.Fatal(err)
		}
	}
	return role
}

func (e *mutationTestEnv) assign(t *testing.T, user *model.User, roles ...*model.Role) {
	t.Helper()
	ids := make([]string, 0, len(roles))
	for _, role := range roles {
		ids = append(ids, role.ID.String())
	}
	if err := e.policyManager.ReplaceRolesForUser(user.ID.String(), ids, authorization.OrgPlatform); err != nil {
		t.Fatal(err)
	}
}

func actorContext(actor *model.User) context.Context {
	return authorization.WithViewer(context.Background(), actor)
}

func TestReplaceUserRolesAcceptsEmptySet(t *testing.T) {
	env := newMutationTestEnv(t, nil)
	actor := env.seedUser(t, "actor")
	target := env.seedUser(t, "target")
	manager := env.seedRole(t, "manager", "role:manage")
	reader := env.seedRole(t, "reader", "comic:read")
	env.assign(t, actor, manager)
	env.assign(t, target, reader)

	result := env.service.ReplaceUserRoles(actorContext(actor), target.ID, []uuid.UUID{}, "g1:u1")
	if !result.Success {
		t.Fatalf("expected success, got %#v", result)
	}
	roles, err := env.policyManager.RolesForUser(target.ID.String(), authorization.OrgPlatform)
	if err != nil || len(roles) != 0 {
		t.Fatalf("expected no roles, got %v, %v", roles, err)
	}
	_, userVersion, err := env.revisions.Current(context.Background(), target.ID)
	if err != nil || userVersion != 2 {
		t.Fatalf("expected user revision 2, got %d, %v", userVersion, err)
	}
	entries, total, err := env.auditRepo.List(context.Background(), authorizationaudit.ListInput{Action: "user.roles_replaced"})
	if err != nil || total != 1 || entries[0].TargetID != target.ID {
		t.Fatalf("unexpected audit: %#v total=%d err=%v", entries, total, err)
	}
	if len(env.cache.deleteKeys) != 1 {
		t.Fatalf("expected old target cache key deletion, got %v", env.cache.deleteKeys)
	}
}

func TestReplaceRolePermissionsAcceptsEmptySet(t *testing.T) {
	env := newMutationTestEnv(t, nil)
	actor := env.seedUser(t, "actor")
	manager := env.seedRole(t, "manager", "role:manage")
	targetRole := env.seedRole(t, "reader", "comic:read")
	env.assign(t, actor, manager)

	result := env.service.ReplaceRolePermissions(actorContext(actor), targetRole.ID, []string{}, "g1")
	if !result.Success {
		t.Fatalf("expected success, got %#v", result)
	}
	permissions, err := env.policyManager.PermissionNamesForRole(targetRole.ID.String(), authorization.OrgPlatform)
	if err != nil || len(permissions) != 0 {
		t.Fatalf("expected no permissions, got %v, %v", permissions, err)
	}
	global, err := env.revisions.CurrentGlobal(context.Background())
	if err != nil || global != 2 {
		t.Fatalf("expected global revision 2, got %d, %v", global, err)
	}
}

func TestReplaceUserRolesRejectsSelfLockout(t *testing.T) {
	env := newMutationTestEnv(t, nil)
	actor := env.seedUser(t, "actor")
	manager := env.seedRole(t, "manager", "role:manage")
	env.assign(t, actor, manager)

	result := env.service.ReplaceUserRoles(actorContext(actor), actor.ID, []uuid.UUID{}, "g1:u1")
	assertConflictCode(t, result, "SELF_MANAGE_REQUIRED")
	roles, err := env.policyManager.RolesForUser(actor.ID.String(), authorization.OrgPlatform)
	if err != nil || len(roles) != 1 || roles[0] != manager.ID.String() {
		t.Fatalf("expected original role to be restored, got %v, %v", roles, err)
	}
}

func TestReplaceRolePermissionsRejectsLastRoleManager(t *testing.T) {
	env := newMutationTestEnv(t, nil)
	actor := env.seedUser(t, "actor-with-stale-request")
	managerUser := env.seedUser(t, "manager")
	managerRole := env.seedRole(t, "manager", "role:manage")
	env.assign(t, managerUser, managerRole)

	result := env.service.ReplaceRolePermissions(actorContext(actor), managerRole.ID, []string{}, "g1")
	assertConflictCode(t, result, "LAST_ROLE_MANAGER")
	permissions, err := env.policyManager.PermissionNamesForRole(managerRole.ID.String(), authorization.OrgPlatform)
	if err != nil || len(permissions) != 1 || permissions[0] != "role:manage" {
		t.Fatalf("expected original permission to be restored, got %v, %v", permissions, err)
	}
}

func TestReplacementRejectsStaleAuthorizationVersion(t *testing.T) {
	env := newMutationTestEnv(t, nil)
	actor := env.seedUser(t, "actor")
	target := env.seedUser(t, "target")
	manager := env.seedRole(t, "manager", "role:manage")
	reader := env.seedRole(t, "reader", "comic:read")
	env.assign(t, actor, manager)
	env.assign(t, target, reader)

	result := env.service.ReplaceUserRoles(actorContext(actor), target.ID, []uuid.UUID{}, "g999:u1")
	assertConflictCode(t, result, "AUTHORIZATION_STATE_CHANGED")
	roles, err := env.policyManager.RolesForUser(target.ID.String(), authorization.OrgPlatform)
	if err != nil || len(roles) != 1 {
		t.Fatalf("stale request changed roles: %v, %v", roles, err)
	}
}

func TestAuditFailureRestoresPreviousCasbinPolicy(t *testing.T) {
	baseDB := testutil.NewSQLiteDB(t)
	_ = baseDB
	env := newMutationTestEnv(t, &failingAuditStore{err: errors.New("audit unavailable")})
	actor := env.seedUser(t, "actor")
	target := env.seedUser(t, "target")
	manager := env.seedRole(t, "manager", "role:manage")
	reader := env.seedRole(t, "reader", "comic:read")
	env.assign(t, actor, manager)
	env.assign(t, target, reader)
	store := env.service.auditRepo.(*failingAuditStore)
	store.base = env.auditRepo

	result := env.service.ReplaceUserRoles(actorContext(actor), target.ID, []uuid.UUID{}, "g1:u1")
	if result.HttpStatus != http.StatusInternalServerError {
		t.Fatalf("expected internal error, got %#v", result)
	}
	roles, err := env.policyManager.RolesForUser(target.ID.String(), authorization.OrgPlatform)
	if err != nil || len(roles) != 1 || roles[0] != reader.ID.String() {
		t.Fatalf("expected policy rollback, got %v, %v", roles, err)
	}
	_, userVersion, err := env.revisions.Current(context.Background(), target.ID)
	if err != nil || userVersion != 1 {
		t.Fatalf("failed mutation bumped revision: %d, %v", userVersion, err)
	}
}

func assertConflictCode(t *testing.T, result response.Result, code string) {
	t.Helper()
	if result.HttpStatus != http.StatusConflict || result.Code != code {
		t.Fatalf("expected conflict %s, got %#v", code, result)
	}
}
