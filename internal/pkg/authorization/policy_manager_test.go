package authorization

import (
	"context"
	"os"
	"testing"

	casbinlib "github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/google/uuid"

	casbinpkg "manga-go/internal/pkg/casbin"
)

func newTestPolicyManager(t *testing.T) (*PolicyManager, *Authorizer) {
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
	return NewPolicyManager(PolicyManagerParams{Enforcer: enforcer}), NewAuthorizer(enforcer)
}

func TestPolicyManagerGrantsReachTheUserThroughTheRole(t *testing.T) {
	pm, authorizer := newTestPolicyManager(t)
	ctx := context.Background()

	userID := uuid.New().String()
	roleID := uuid.New().String()

	if err := pm.AddRoleForUser(userID, roleID, OrgPlatform); err != nil {
		t.Fatalf("failed to add role for user: %v", err)
	}
	if err := pm.ReplacePermissionsForRole(roleID, []string{"comic:read"}, OrgPlatform); err != nil {
		t.Fatalf("failed to grant permission to role: %v", err)
	}

	err := authorizer.Enforce(ctx, Request{
		Subject: userID,
		Org:     OrgPlatform,
		Action:  ActionRead,
		Object:  ObjectComic,
		Context: CtxAny,
	})
	if err != nil {
		t.Fatalf("expected user to be allowed through roleID -> permissionID, got: %v", err)
	}

	err = authorizer.Enforce(ctx, Request{
		Subject: userID,
		Org:     OrgPlatform,
		Action:  ActionRead,
		Object:  ObjectChapter,
		Context: CtxAny,
	})
	if err == nil {
		t.Fatalf("expected unrelated permission to be denied")
	}
}

// The baseline is what a signed-in user gets with no role at all. It used to be
// a switch statement in Go; it is now policy, so it has to be written first.
func TestBaselinePoliciesGiveSignedInUsersTheBasics(t *testing.T) {
	pm, authorizer := newTestPolicyManager(t)
	ctx := context.Background()

	if err := pm.ReplaceBaselinePolicies(); err != nil {
		t.Fatalf("failed to write the baseline policies: %v", err)
	}

	userID := uuid.New().String()
	err := authorizer.Enforce(ctx, Request{
		Subject: userID,
		Org:     OrgPlatform,
		Action:  ActionCreate,
		Object:  ObjectComment,
		Context: CtxAny,
	})
	if err != nil {
		t.Fatalf("expected a signed-in user to create a comment from the baseline, got: %v", err)
	}

	err = authorizer.Enforce(ctx, Request{
		Subject: userID,
		Org:     OrgPlatform,
		Action:  ActionManage,
		Object:  ObjectRole,
		Context: CtxAny,
	})
	if err == nil {
		t.Fatalf("expected implicit reader baseline to deny role management")
	}

	err = authorizer.Enforce(ctx, Request{
		Subject: userID,
		Org:     TranslationGroupOrg(uuid.New()),
		Action:  ActionRead,
		Object:  ObjectChapter,
		Context: CtxPublished,
	})
	if err != nil {
		t.Fatalf("expected authenticated user to read published chapter in group org, got: %v", err)
	}
}

func TestPolicyManagerKeepsTextForTranslationGroupRoles(t *testing.T) {
	pm, authorizer := newTestPolicyManager(t)
	ctx := context.Background()

	userID := uuid.New().String()
	groupID := uuid.New().String()

	if err := pm.AddTranslationGroupMember(userID, groupID); err != nil {
		t.Fatalf("failed to add group member policy: %v", err)
	}

	err := authorizer.Enforce(ctx, Request{
		Subject: userID,
		Org:     TranslationGroupOrgString(groupID),
		Action:  ActionCreate,
		Object:  ObjectChapter,
		Context: CtxGroupMember,
	})
	if err != nil {
		t.Fatalf("expected group member to create chapter through text policies, got: %v", err)
	}

	// Roles are stored under the namespaced org ("tg:<id>"), not the bare group id.
	domain := string(TranslationGroupOrgString(groupID))
	if roles := pm.enforcer.GetRolesForUserInDomain(userID, domain); len(roles) != 1 || roles[0] != roleGroupMember {
		t.Fatalf("expected group member role text, got %#v", roles)
	}
}

func TestPolicyManagerAllowsGroupOwnerToManageComicsInOwnGroup(t *testing.T) {
	pm, authorizer := newTestPolicyManager(t)
	ctx := context.Background()

	userID := uuid.New().String()
	groupID := uuid.New().String()
	otherGroupID := uuid.New()

	if err := pm.AddTranslationGroupOwner(userID, groupID); err != nil {
		t.Fatalf("failed to add group owner policy: %v", err)
	}

	err := authorizer.Enforce(ctx, Request{
		Subject: userID,
		Org:     TranslationGroupOrgString(groupID),
		Action:  ActionDelete,
		Object:  ObjectComic,
		Context: CtxGroupMember,
	})
	if err != nil {
		t.Fatalf("expected group owner to manage comic in own group, got: %v", err)
	}

	err = authorizer.Enforce(ctx, Request{
		Subject: userID,
		Org:     TranslationGroupOrg(otherGroupID),
		Action:  ActionDelete,
		Object:  ObjectComic,
		Context: CtxGroupMember,
	})
	if err == nil {
		t.Fatalf("expected group owner to be denied outside own group")
	}
}
