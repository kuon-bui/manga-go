package seeder

import (
	"context"
	"os"
	"testing"

	"manga-go/internal/pkg/authorization"
	casbinpkg "manga-go/internal/pkg/casbin"
	"manga-go/internal/pkg/config"
	rolerepo "manga-go/internal/pkg/repo/role"
	userrepo "manga-go/internal/pkg/repo/user"
	roleseeder "manga-go/internal/pkg/seeder/role"
	userseeder "manga-go/internal/pkg/seeder/user"
	"manga-go/internal/pkg/testutil"

	casbinlib "github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormlogger "gorm.io/gorm/logger"
)

const (
	seededAdminEmail      = "admin@manga.com"
	seededTranslatorEmail = "seed-user-01@manga.local"
)

type seededAuthz struct {
	db         *gorm.DB
	policy     *authorization.PolicyManager
	authorizer *authorization.Authorizer
}

func newSeededAuthz(t *testing.T) seededAuthz {
	t.Helper()

	data, err := os.ReadFile("../casbin/model.conf")
	if err != nil {
		t.Fatalf("failed to read casbin model: %v", err)
	}
	m, err := model.NewModelFromString(string(data))
	if err != nil {
		t.Fatalf("failed to create casbin model: %v", err)
	}
	e, err := casbinlib.NewEnforcer(m)
	if err != nil {
		t.Fatalf("failed to create casbin enforcer: %v", err)
	}
	enforcer := &casbinpkg.Enforcer{Enforcer: e}

	db := testutil.NewSQLiteDB(t)
	// The seeders probe for existing rows on purpose; GORM logs every miss as an
	// error, which would bury the actual test output.
	db.Logger = gormlogger.Discard
	testutil.MustSyncSchemas(t, db,
		&testutil.Role{},
		&testutil.User{},
	)

	return seededAuthz{
		db:         db,
		policy:     authorization.NewPolicyManager(authorization.PolicyManagerParams{Enforcer: enforcer}),
		authorizer: authorization.NewAuthorizer(enforcer),
	}
}

func (s seededAuthz) runSeeders(t *testing.T) {
	t.Helper()

	f := faker.New()
	cfg := &config.Config{}

	seeders := []Seeder{
		roleseeder.NewRoleSeeder(rolerepo.NewRoleRepo(s.db), f, s.policy),
		userseeder.NewUserSeeder(userrepo.NewUserRepository(s.db, nil), rolerepo.NewRoleRepo(s.db), cfg, f, s.policy),
	}
	for _, seeder := range seeders {
		if err := seeder.Seed(s.db); err != nil {
			t.Fatalf("%s failed: %v", seeder.Name(), err)
		}
	}
}

func (s seededAuthz) userIDByEmail(t *testing.T, email string) uuid.UUID {
	t.Helper()

	user, err := userrepo.NewUserRepository(s.db, nil).FindOne(context.Background(), []any{
		clause.Eq{Column: "email", Value: email},
	}, nil)
	if err != nil {
		t.Fatalf("failed to find seeded user %s: %v", email, err)
	}
	return user.ID
}

func (s seededAuthz) allows(t *testing.T, userID uuid.UUID, action authorization.Action, object authorization.Object) bool {
	t.Helper()
	return s.allowsSubject(t, authorization.Subject(userID), action, object, authorization.CtxAny)
}

func (s seededAuthz) allowsSubject(t *testing.T, subject string, action authorization.Action, object authorization.Object, ctx authorization.Context) bool {
	t.Helper()

	err := s.authorizer.Enforce(context.Background(), authorization.Request{
		Subject: subject,
		Org:     authorization.OrgPlatform,
		Action:  action,
		Object:  object,
		Context: ctx,
	})
	if err != nil && !errorIsForbidden(err) {
		t.Fatalf("enforce failed for an unexpected reason: %v", err)
	}
	return err == nil
}

func errorIsForbidden(err error) bool {
	return err == authorization.ErrForbidden
}

// Seeding roles and users must leave the policy engine able to answer for them:
// a seeded admin that exists only in the database is denied everything.
func TestSeededRolesAreEnforceable(t *testing.T) {
	env := newSeededAuthz(t)
	env.runSeeders(t)

	adminID := env.userIDByEmail(t, seededAdminEmail)
	translatorID := env.userIDByEmail(t, seededTranslatorEmail)

	cases := []struct {
		name    string
		userID  uuid.UUID
		action  authorization.Action
		object  authorization.Object
		allowed bool
	}{
		{"admin manages roles", adminID, authorization.ActionManage, authorization.ObjectRole, true},
		{"admin deletes comics", adminID, authorization.ActionDelete, authorization.ObjectComic, true},
		{"admin manages users", adminID, authorization.ActionManage, authorization.ObjectUser, true},
		// The right /admin/comic-stats requires.
		{"admin updates comics", adminID, authorization.ActionUpdate, authorization.ObjectComic, true},
		{"translator cannot update comics platform-wide", translatorID, authorization.ActionUpdate, authorization.ObjectComic, false},
		{"translator creates chapters", translatorID, authorization.ActionCreate, authorization.ObjectChapter, true},
		{"translator cannot manage roles", translatorID, authorization.ActionManage, authorization.ObjectRole, false},
		{"translator cannot delete comics", translatorID, authorization.ActionDelete, authorization.ObjectComic, false},
		// The role-administration UI needs the catalog to render.
		{"admin reads the permission catalog", adminID, authorization.ActionRead, authorization.ObjectPermission, true},
		{"translator cannot read the permission catalog", translatorID, authorization.ActionRead, authorization.ObjectPermission, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := env.allows(t, tc.userID, tc.action, tc.object); got != tc.allowed {
				t.Fatalf("expected allowed=%v, got %v", tc.allowed, got)
			}
		})
	}
}

// Seeding also has to establish the baseline: what a visitor may do before any
// role is involved. These rights used to be a switch statement in Go.
func TestSeedingEstablishesTheBaseline(t *testing.T) {
	env := newSeededAuthz(t)
	env.runSeeders(t)

	anyUser := authorization.Subject(uuid.New())

	if !env.allowsSubject(t, authorization.SubjectAnonymous, authorization.ActionRead, authorization.ObjectComic, authorization.CtxPublished) {
		t.Error("an anonymous visitor should be able to read a published comic")
	}
	if env.allowsSubject(t, authorization.SubjectAnonymous, authorization.ActionRead, authorization.ObjectComic, authorization.CtxAny) {
		t.Error("an anonymous visitor must not be able to read an unpublished comic")
	}
	if env.allowsSubject(t, authorization.SubjectAnonymous, authorization.ActionCreate, authorization.ObjectComment, authorization.CtxAny) {
		t.Error("an anonymous visitor must not be able to comment")
	}
	if !env.allowsSubject(t, anyUser, authorization.ActionCreate, authorization.ObjectComment, authorization.CtxAny) {
		t.Error("any signed-in user should be able to comment")
	}
	if !env.allowsSubject(t, anyUser, authorization.ActionUpdate, authorization.ObjectUser, authorization.CtxSelf) {
		t.Error("any signed-in user should be able to update their own profile")
	}
	if env.allowsSubject(t, anyUser, authorization.ActionUpdate, authorization.ObjectUser, authorization.CtxAny) {
		t.Error("the self-update right must not extend to other users")
	}
}

// Re-running the seeder is the documented way to repair an environment, so it
// must converge rather than pile up or wipe out policy rules.
func TestReSeedingKeepsAuthorizationIntact(t *testing.T) {
	env := newSeededAuthz(t)
	env.runSeeders(t)
	env.runSeeders(t)

	adminID := env.userIDByEmail(t, seededAdminEmail)
	translatorID := env.userIDByEmail(t, seededTranslatorEmail)

	if !env.allows(t, adminID, authorization.ActionManage, authorization.ObjectRole) {
		t.Error("admin should keep its permissions after a re-seed")
	}
	if env.allows(t, translatorID, authorization.ActionManage, authorization.ObjectRole) {
		t.Error("translator must not gain permissions from a re-seed")
	}
}
