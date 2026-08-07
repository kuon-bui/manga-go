package authorizationadmin

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"
	rolerepo "manga-go/internal/pkg/repo/role"
	"manga-go/internal/pkg/testutil"

	"github.com/google/uuid"
	dto "github.com/prometheus/client_model/go"
	"gorm.io/gorm"
)

type fakeProfileCache struct {
	profile    *AuthorizationProfile
	getErr     error
	setErr     error
	getKey     string
	setKey     string
	setTTL     time.Duration
	setCalls   int
	deleteKeys []string
}

func (f *fakeProfileCache) Get(_ context.Context, key string, target *AuthorizationProfile) (bool, error) {
	f.getKey = key
	if f.getErr != nil {
		return false, f.getErr
	}
	if f.profile == nil {
		return false, nil
	}
	*target = *f.profile
	return true, nil
}

func (f *fakeProfileCache) Set(_ context.Context, key string, profile *AuthorizationProfile, ttl time.Duration) error {
	f.setCalls++
	f.setKey = key
	f.setTTL = ttl
	f.profile = profile
	return f.setErr
}

func (f *fakeProfileCache) Delete(_ context.Context, key string) error {
	f.deleteKeys = append(f.deleteKeys, key)
	return nil
}

type fakeRevisionStore struct {
	global uint64
	user   uint64
	err    error
}

func (f *fakeRevisionStore) Current(context.Context, uuid.UUID) (uint64, uint64, error) {
	return f.global, f.user, f.err
}

func (f *fakeRevisionStore) CurrentGlobal(context.Context) (uint64, error) {
	return f.global, f.err
}

func (f *fakeRevisionStore) CurrentMany(_ context.Context, userIDs []uuid.UUID) (uint64, map[uuid.UUID]uint64, error) {
	users := make(map[uuid.UUID]uint64, len(userIDs))
	for _, userID := range userIDs {
		users[userID] = f.user
	}
	return f.global, users, f.err
}

func (f *fakeRevisionStore) BumpGlobalTx(*gorm.DB) (uint64, error) {
	f.global++
	return f.global, nil
}

func (f *fakeRevisionStore) BumpUserTx(*gorm.DB, uuid.UUID) (uint64, error) {
	f.user++
	return f.user, nil
}

type profileTestEnv struct {
	service       *Service
	cache         *fakeProfileCache
	revisions     *fakeRevisionStore
	policyManager *authorization.PolicyManager
	enforcer      interface {
		AddPolicy(...any) (bool, error)
	}
	db *gorm.DB
}

func newProfileTestEnv(t *testing.T) *profileTestEnv {
	t.Helper()

	db := testutil.NewSQLiteDB(t)
	testutil.MustSyncSchemas(t, db, &testutil.Role{})
	enforcer := testutil.NewInMemoryEnforcer(t)
	policyManager := authorization.NewPolicyManager(authorization.PolicyManagerParams{Enforcer: enforcer})
	cache := &fakeProfileCache{}
	revisions := &fakeRevisionStore{global: 1, user: 1}

	return &profileTestEnv{
		service: &Service{
			logger:        logger.NewLogger(),
			roleRepo:      rolerepo.NewRoleRepo(db),
			policyManager: policyManager,
			authorizer:    authorization.NewAuthorizer(enforcer),
			revisions:     revisions,
			cache:         cache,
		},
		cache:         cache,
		revisions:     revisions,
		policyManager: policyManager,
		enforcer:      enforcer,
		db:            db,
	}
}

func (e *profileTestEnv) seedUserAndRole(t *testing.T, roleName string) (uuid.UUID, uuid.UUID) {
	t.Helper()

	userID := uuid.New()
	roleID := uuid.New()
	if err := e.db.Exec(
		"INSERT INTO roles (id, name, created_at, updated_at) VALUES (?, ?, datetime('now'), datetime('now'))",
		roleID, roleName,
	).Error; err != nil {
		t.Fatal(err)
	}
	if err := e.policyManager.AddRoleForUser(userID.String(), roleID.String(), authorization.OrgPlatform); err != nil {
		t.Fatal(err)
	}
	return userID, roleID
}

func (e *profileTestEnv) mustAssign(t *testing.T, roleID uuid.UUID, permissions []string) {
	t.Helper()
	if err := e.policyManager.ReplacePermissionsForRole(roleID.String(), permissions, authorization.OrgPlatform); err != nil {
		t.Fatal(err)
	}
}

func TestGetProfileExpandsManageIntoEffectiveCatalogNames(t *testing.T) {
	env := newProfileTestEnv(t)
	userID, roleID := env.seedUserAndRole(t, "manager")
	env.mustAssign(t, roleID, []string{"user:manage", "role:manage"})

	profile, err := env.service.GetProfile(context.Background(), userID)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"user:read", "user:write", "user:delete", "user:manage", "role:manage"} {
		if !slices.Contains(profile.Permissions, name) {
			t.Fatalf("expected %s in %v", name, profile.Permissions)
		}
	}
	if profile.Version != "g1:u1" || env.cache.setTTL != 10*time.Minute {
		t.Fatalf("unexpected version or TTL: %q, %s", profile.Version, env.cache.setTTL)
	}
}

func TestGetProfileRequiresEveryExpandedWriteGrant(t *testing.T) {
	env := newProfileTestEnv(t)
	userID, roleID := env.seedUserAndRole(t, "partial-writer")

	for _, action := range []authorization.Action{authorization.ActionCreate, authorization.ActionUpdate} {
		if _, err := env.enforcer.AddPolicy(
			string(authorization.OrgPlatform), roleID.String(), string(action),
			string(authorization.ObjectComic), string(authorization.CtxAny), "allow",
		); err != nil {
			t.Fatal(err)
		}
	}

	profile, err := env.service.GetProfile(context.Background(), userID)
	if err != nil {
		t.Fatal(err)
	}
	if slices.Contains(profile.Permissions, "comic:write") {
		t.Fatalf("partial concrete grants must not expose comic:write: %v", profile.Permissions)
	}

	env.cache.profile = nil
	if _, err := env.enforcer.AddPolicy(
		string(authorization.OrgPlatform), roleID.String(), string(authorization.ActionPublish),
		string(authorization.ObjectComic), string(authorization.CtxAny), "allow",
	); err != nil {
		t.Fatal(err)
	}
	profile, err = env.service.GetProfile(context.Background(), userID)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(profile.Permissions, "comic:write") {
		t.Fatalf("all concrete grants must expose comic:write: %v", profile.Permissions)
	}
}

func TestGetProfileFallsBackWhenRedisIsUnavailable(t *testing.T) {
	env := newProfileTestEnv(t)
	env.cache.getErr = errors.New("redis unavailable")
	userID, roleID := env.seedUserAndRole(t, "reader")
	env.mustAssign(t, roleID, []string{"comic:read"})

	profile, err := env.service.GetProfile(context.Background(), userID)
	if err != nil || profile.UserID != userID {
		t.Fatalf("expected live Casbin profile, got %#v, %v", profile, err)
	}
	if env.cache.setCalls != 0 {
		t.Fatal("must not write cache after a Redis read failure")
	}
}

func TestGetProfileMeasuresRedisReadFailure(t *testing.T) {
	env := newProfileTestEnv(t)
	env.cache.getErr = errors.New("redis unavailable")
	userID, roleID := env.seedUserAndRole(t, "reader")
	env.mustAssign(t, roleID, []string{"comic:read"})
	before := counterValue(t, authorizationProfileCacheFailures.WithLabelValues("read"))

	if _, err := env.service.GetProfile(context.Background(), userID); err != nil {
		t.Fatal(err)
	}
	after := counterValue(t, authorizationProfileCacheFailures.WithLabelValues("read"))
	if after != before+1 {
		t.Fatalf("expected one measured read failure, before=%v after=%v", before, after)
	}
}

func counterValue(t *testing.T, counter interface{ Write(*dto.Metric) error }) float64 {
	t.Helper()
	metric := &dto.Metric{}
	if err := counter.Write(metric); err != nil {
		t.Fatal(err)
	}
	return metric.GetCounter().GetValue()
}

func TestGetProfileReturnsLiveDataWhenRedisWriteFails(t *testing.T) {
	env := newProfileTestEnv(t)
	env.cache.setErr = errors.New("redis write unavailable")
	userID, roleID := env.seedUserAndRole(t, "reader")
	env.mustAssign(t, roleID, []string{"comic:read"})

	profile, err := env.service.GetProfile(context.Background(), userID)
	if err != nil || !slices.Contains(profile.Permissions, "comic:read") {
		t.Fatalf("expected live profile despite cache write failure, got %#v, %v", profile, err)
	}
}

func TestGetProfileCacheHitAvoidsLiveDependencies(t *testing.T) {
	userID := uuid.MustParse("4d8ca2fe-740f-4b75-8776-6e827cb89540")
	cache := &fakeProfileCache{profile: &AuthorizationProfile{UserID: userID, Version: "g1:u1"}}
	service := &Service{
		cache:     cache,
		revisions: &fakeRevisionStore{global: 1, user: 1},
	}

	profile, err := service.GetProfile(context.Background(), userID)
	if err != nil || profile.UserID != userID {
		t.Fatalf("expected cached profile, got %#v, %v", profile, err)
	}
	wantKey := "authorization:profile:g1:u1:user:4d8ca2fe-740f-4b75-8776-6e827cb89540"
	if cache.getKey != wantKey || cache.setCalls != 0 {
		t.Fatalf("unexpected cache behavior: key=%q setCalls=%d", cache.getKey, cache.setCalls)
	}
}

func TestEffectivePermissionPropagatesAuthorizerFailure(t *testing.T) {
	service := &Service{authorizer: authorization.NewAuthorizer(nil)}
	definition, ok := authorization.LookupPermission("comic:read")
	if !ok {
		t.Fatal("missing catalog definition")
	}

	allowed, err := service.isEffective(context.Background(), uuid.NewString(), definition)
	if allowed || !errors.Is(err, authorization.ErrAuthorizerUnavailable) {
		t.Fatalf("expected authorizer failure, got allowed=%v err=%v", allowed, err)
	}
}
