package authorization

import (
	"os"
	"testing"

	casbinpkg "manga-go/internal/pkg/casbin"

	casbinlib "github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/google/uuid"
)

// These tests exercise model.conf itself, so they call the enforcer directly.
// Going through Authorizer would mix in whatever shortcuts that layer applies.
func newModelTestAuthorizer(t *testing.T) (*casbinpkg.Enforcer, *casbinpkg.Enforcer) {
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
	return enforcer, enforcer
}

func mustAddPolicy(t *testing.T, e *casbinpkg.Enforcer, org Org, subject string, action Action, object Object, ctx Context, effect string) {
	t.Helper()

	if _, err := e.AddPolicy(string(org), subject, string(action), string(object), string(ctx), effect); err != nil {
		t.Fatalf("failed to add policy: %v", err)
	}
}

func allowed(t *testing.T, e *casbinpkg.Enforcer, subject string, org Org, action Action, object Object, ctx Context) bool {
	t.Helper()

	ok, err := e.Enforce(subject, string(org), string(action), string(object), string(ctx))
	if err != nil {
		t.Fatalf("enforce failed: %v", err)
	}
	return ok
}

// A policy granted to "anonymous" applies to everyone, signed in or not: it is
// the floor of what the platform is willing to show.
func TestAnonymousPoliciesApplyWithoutAnyRole(t *testing.T) {
	e, a := newModelTestAuthorizer(t)
	mustAddPolicy(t, e, OrgPlatform, SubjectAnonymous, ActionRead, ObjectComic, CtxPublished, "allow")

	if !allowed(t, a, SubjectAnonymous, OrgPlatform, ActionRead, ObjectComic, CtxPublished) {
		t.Error("an anonymous visitor should be able to read a published comic")
	}
	if !allowed(t, a, uuid.NewString(), OrgPlatform, ActionRead, ObjectComic, CtxPublished) {
		t.Error("a signed-in user should also get anonymous-level access")
	}
}

// "authenticated" is the complement: signed-in users get it implicitly, without
// a grouping row per user, and anonymous visitors never do.
func TestAuthenticatedPoliciesRequireASignedInSubject(t *testing.T) {
	e, a := newModelTestAuthorizer(t)
	mustAddPolicy(t, e, OrgPlatform, RoleAuthenticated, ActionCreate, ObjectComment, CtxAny, "allow")

	if !allowed(t, a, uuid.NewString(), OrgPlatform, ActionCreate, ObjectComment, CtxAny) {
		t.Error("a signed-in user should be able to create a comment")
	}
	if allowed(t, a, SubjectAnonymous, OrgPlatform, ActionCreate, ObjectComment, CtxAny) {
		t.Error("an anonymous visitor must not be able to create a comment")
	}
}

func TestRolePoliciesApplyThroughASingleGroupingLevel(t *testing.T) {
	e, a := newModelTestAuthorizer(t)
	userID := uuid.NewString()
	roleID := uuid.NewString()

	mustAddPolicy(t, e, OrgPlatform, roleID, ActionDelete, ObjectComic, CtxAny, "allow")
	if _, err := e.AddGroupingPolicy(userID, roleID, string(OrgPlatform)); err != nil {
		t.Fatalf("failed to add grouping policy: %v", err)
	}

	if !allowed(t, a, userID, OrgPlatform, ActionDelete, ObjectComic, CtxAny) {
		t.Error("a user holding the role should inherit its policy")
	}
	if allowed(t, a, uuid.NewString(), OrgPlatform, ActionDelete, ObjectComic, CtxAny) {
		t.Error("a user without the role must be denied")
	}
}

func TestManageActionCoversTheOtherActions(t *testing.T) {
	e, a := newModelTestAuthorizer(t)
	userID := uuid.NewString()
	roleID := uuid.NewString()

	mustAddPolicy(t, e, OrgPlatform, roleID, ActionManage, ObjectComic, CtxAny, "allow")
	if _, err := e.AddGroupingPolicy(userID, roleID, string(OrgPlatform)); err != nil {
		t.Fatalf("failed to add grouping policy: %v", err)
	}

	for _, action := range []Action{ActionRead, ActionCreate, ActionUpdate, ActionDelete, ActionPublish} {
		if !allowed(t, a, userID, OrgPlatform, action, ObjectComic, CtxAny) {
			t.Errorf("manage should cover %q", action)
		}
	}
}

func TestPlatformPoliciesReachEveryOrgButOrgPoliciesStayPut(t *testing.T) {
	e, a := newModelTestAuthorizer(t)
	userID := uuid.NewString()
	roleID := uuid.NewString()
	groupA := TranslationGroupOrg(uuid.New())
	groupB := TranslationGroupOrg(uuid.New())

	mustAddPolicy(t, e, groupA, roleID, ActionCreate, ObjectChapter, CtxGroupMember, "allow")
	if _, err := e.AddGroupingPolicy(userID, roleID, string(groupA)); err != nil {
		t.Fatalf("failed to add grouping policy: %v", err)
	}

	if !allowed(t, a, userID, groupA, ActionCreate, ObjectChapter, CtxGroupMember) {
		t.Error("the group policy should apply inside its own group")
	}
	if allowed(t, a, userID, groupB, ActionCreate, ObjectChapter, CtxGroupMember) {
		t.Error("a group policy must not leak into another group")
	}
}

// The policy table declares an effect column. If deny cannot actually override
// allow, that column is a trap: someone will write a deny rule and believe it.
func TestDenyOverridesAllow(t *testing.T) {
	e, a := newModelTestAuthorizer(t)
	userID := uuid.NewString()
	roleID := uuid.NewString()

	mustAddPolicy(t, e, OrgPlatform, roleID, ActionManage, ObjectComic, CtxAny, "allow")
	mustAddPolicy(t, e, OrgPlatform, roleID, ActionDelete, ObjectComic, CtxAny, "deny")
	if _, err := e.AddGroupingPolicy(userID, roleID, string(OrgPlatform)); err != nil {
		t.Fatalf("failed to add grouping policy: %v", err)
	}

	if !allowed(t, a, userID, OrgPlatform, ActionUpdate, ObjectComic, CtxAny) {
		t.Error("the allow rule should still apply to actions that are not denied")
	}
	if allowed(t, a, userID, OrgPlatform, ActionDelete, ObjectComic, CtxAny) {
		t.Error("an explicit deny must override a broader allow")
	}
}

func TestContextMustMatchUnlessThePolicySaysAny(t *testing.T) {
	e, a := newModelTestAuthorizer(t)
	userID := uuid.NewString()
	roleID := uuid.NewString()

	mustAddPolicy(t, e, OrgPlatform, roleID, ActionUpdate, ObjectComment, CtxOwner, "allow")
	if _, err := e.AddGroupingPolicy(userID, roleID, string(OrgPlatform)); err != nil {
		t.Fatalf("failed to add grouping policy: %v", err)
	}

	if !allowed(t, a, userID, OrgPlatform, ActionUpdate, ObjectComment, CtxOwner) {
		t.Error("the owner context should match")
	}
	if allowed(t, a, userID, OrgPlatform, ActionUpdate, ObjectComment, CtxAny) {
		t.Error("an owner-scoped policy must not satisfy an unscoped request")
	}
}
