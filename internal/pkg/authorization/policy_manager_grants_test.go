package authorization

import (
	"context"
	"sort"
	"testing"

	"github.com/google/uuid"
)

func allowedVia(t *testing.T, a *Authorizer, subject string, org Org, action Action, object Object, ctx Context) bool {
	t.Helper()

	err := a.Enforce(context.Background(), Request{
		Subject: subject,
		Org:     org,
		Action:  action,
		Object:  object,
		Context: ctx,
	})
	return err == nil
}

func grantedNames(t *testing.T, pm *PolicyManager, roleID string, org Org) []string {
	t.Helper()

	names, err := pm.PermissionNamesForRole(roleID, org)
	if err != nil {
		t.Fatalf("failed to read back the role's permissions: %v", err)
	}
	sort.Strings(names)
	return names
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// A role grant is written straight onto the role: no intermediate permission
// entity, so nothing can be left dangling when a permission is renamed.
func TestReplacePermissionsForRoleWritesPolicyAgainstTheRole(t *testing.T) {
	pm, _ := newTestPolicyManager(t)
	roleID := uuid.NewString()

	if err := pm.ReplacePermissionsForRole(roleID, []string{"comic:read", "role:manage"}, OrgPlatform); err != nil {
		t.Fatalf("failed to grant permissions: %v", err)
	}

	if got, want := grantedNames(t, pm, roleID, OrgPlatform), []string{"comic:read", "role:manage"}; !equalStrings(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestReplacePermissionsForRoleRoundTripsWriteShorthand(t *testing.T) {
	pm, _ := newTestPolicyManager(t)
	roleID := uuid.NewString()

	if err := pm.ReplacePermissionsForRole(roleID, []string{"comic:write"}, OrgPlatform); err != nil {
		t.Fatalf("failed to grant permissions: %v", err)
	}

	if got, want := grantedNames(t, pm, roleID, OrgPlatform), []string{"comic:write"}; !equalStrings(got, want) {
		t.Fatalf("expected the create/update/publish rules to read back as comic:write, got %v", got)
	}
}

func TestReplacePermissionsForRoleReplacesRatherThanAccumulates(t *testing.T) {
	pm, _ := newTestPolicyManager(t)
	roleID := uuid.NewString()

	if err := pm.ReplacePermissionsForRole(roleID, []string{"comic:read", "comic:delete"}, OrgPlatform); err != nil {
		t.Fatalf("failed to grant permissions: %v", err)
	}
	if err := pm.ReplacePermissionsForRole(roleID, []string{"comic:read"}, OrgPlatform); err != nil {
		t.Fatalf("failed to replace permissions: %v", err)
	}

	if got, want := grantedNames(t, pm, roleID, OrgPlatform), []string{"comic:read"}; !equalStrings(got, want) {
		t.Fatalf("expected the replaced set to be exactly comic:read, got %v", got)
	}
}

func TestReplacePermissionsForRoleRejectsUnknownNameBeforeWriting(t *testing.T) {
	pm, _ := newTestPolicyManager(t)
	roleID := uuid.NewString()

	if err := pm.ReplacePermissionsForRole(roleID, []string{"comic:read"}, OrgPlatform); err != nil {
		t.Fatalf("failed to grant permissions: %v", err)
	}

	err := pm.ReplacePermissionsForRole(roleID, []string{"comic:delete", "comic:frobnicate"}, OrgPlatform)
	if err == nil {
		t.Fatal("expected an unknown permission name to be rejected")
	}

	// The rejected call must not have stripped what the role already had.
	if got, want := grantedNames(t, pm, roleID, OrgPlatform), []string{"comic:read"}; !equalStrings(got, want) {
		t.Fatalf("a rejected grant must leave the role untouched, got %v", got)
	}
}

func TestRevokePermissionFromRoleRemovesOnlyThatPermission(t *testing.T) {
	pm, _ := newTestPolicyManager(t)
	roleID := uuid.NewString()

	if err := pm.ReplacePermissionsForRole(roleID, []string{"comic:read", "comic:delete", "role:manage"}, OrgPlatform); err != nil {
		t.Fatalf("failed to grant permissions: %v", err)
	}
	if err := pm.RevokePermissionFromRole(roleID, "comic:delete", OrgPlatform); err != nil {
		t.Fatalf("failed to revoke: %v", err)
	}

	if got, want := grantedNames(t, pm, roleID, OrgPlatform), []string{"comic:read", "role:manage"}; !equalStrings(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestGrantPermissionToRoleAddsWithoutDisturbingTheRest(t *testing.T) {
	pm, _ := newTestPolicyManager(t)
	roleID := uuid.NewString()

	if err := pm.ReplacePermissionsForRole(roleID, []string{"comic:read"}, OrgPlatform); err != nil {
		t.Fatalf("failed to grant permissions: %v", err)
	}
	if err := pm.GrantPermissionToRole(roleID, "role:manage", OrgPlatform); err != nil {
		t.Fatalf("failed to grant: %v", err)
	}

	if got, want := grantedNames(t, pm, roleID, OrgPlatform), []string{"comic:read", "role:manage"}; !equalStrings(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// Deleting a role must take its grants with it; otherwise the rules linger
// under a role id nothing points at any more.
func TestRemoveRoleClearsItsGrantsAndItsHolders(t *testing.T) {
	pm, _ := newTestPolicyManager(t)
	roleID := uuid.NewString()
	userID := uuid.NewString()

	if err := pm.ReplacePermissionsForRole(roleID, []string{"comic:read"}, OrgPlatform); err != nil {
		t.Fatalf("failed to grant permissions: %v", err)
	}
	if err := pm.AddRoleForUser(userID, roleID, OrgPlatform); err != nil {
		t.Fatalf("failed to assign role: %v", err)
	}

	if err := pm.RemoveRole(roleID, OrgPlatform); err != nil {
		t.Fatalf("failed to remove role: %v", err)
	}

	if got := grantedNames(t, pm, roleID, OrgPlatform); len(got) != 0 {
		t.Errorf("expected the role's grants to be gone, got %v", got)
	}
	roles, err := pm.RolesForUser(userID, OrgPlatform)
	if err != nil {
		t.Fatalf("failed to read the user's roles: %v", err)
	}
	if len(roles) != 0 {
		t.Errorf("expected the user to no longer hold the role, got %v", roles)
	}
}

func TestRolesForUserReportsAssignedRoles(t *testing.T) {
	pm, _ := newTestPolicyManager(t)
	userID := uuid.NewString()
	first := uuid.NewString()
	second := uuid.NewString()

	if err := pm.ReplaceRolesForUser(userID, []string{first, second}, OrgPlatform); err != nil {
		t.Fatalf("failed to assign roles: %v", err)
	}

	roles, err := pm.RolesForUser(userID, OrgPlatform)
	if err != nil {
		t.Fatalf("failed to read the user's roles: %v", err)
	}
	sort.Strings(roles)
	want := []string{first, second}
	sort.Strings(want)
	if !equalStrings(roles, want) {
		t.Fatalf("expected %v, got %v", want, roles)
	}
}

// The baseline is what every visitor gets before any role is involved.
func TestReplaceBaselinePoliciesIsIdempotent(t *testing.T) {
	pm, authorizer := newTestPolicyManager(t)

	for range 2 {
		if err := pm.ReplaceBaselinePolicies(); err != nil {
			t.Fatalf("failed to write the baseline policies: %v", err)
		}
	}

	if !allowedVia(t, authorizer, SubjectAnonymous, OrgPlatform, ActionRead, ObjectComic, CtxPublished) {
		t.Error("an anonymous visitor should be able to read a published comic")
	}
	if !allowedVia(t, authorizer, uuid.NewString(), OrgPlatform, ActionCreate, ObjectComment, CtxAny) {
		t.Error("a signed-in user should be able to comment")
	}
	if allowedVia(t, authorizer, SubjectAnonymous, OrgPlatform, ActionCreate, ObjectComment, CtxAny) {
		t.Error("an anonymous visitor must not be able to comment")
	}
	if allowedVia(t, authorizer, uuid.NewString(), OrgPlatform, ActionDelete, ObjectComic, CtxAny) {
		t.Error("the baseline must not grant comic deletion to everyone")
	}
}
