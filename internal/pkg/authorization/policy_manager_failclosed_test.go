package authorization

import (
	"testing"

	"github.com/google/uuid"
)

func TestPolicyManagerRejectsInvalidPermissionName(t *testing.T) {
	pm, _ := newTestPolicyManager(t)

	// A dot instead of a colon, and an action the catalog does not define.
	for _, name := range []string{"comics.read", "comic:frobnicate"} {
		if err := pm.GrantPermissionToRole(uuid.NewString(), name, OrgPlatform); err == nil {
			t.Errorf("expected %q to be rejected, got silent success", name)
		}
	}
}

func TestPolicyManagerLeavesNoPolicyBehindForInvalidPermissionName(t *testing.T) {
	pm, _ := newTestPolicyManager(t)
	roleID := uuid.NewString()

	_ = pm.GrantPermissionToRole(roleID, "comics.read", OrgPlatform)

	names, err := pm.PermissionNamesForRole(roleID, OrgPlatform)
	if err != nil {
		t.Fatalf("failed to read the role's permissions: %v", err)
	}
	if len(names) != 0 {
		t.Fatalf("expected no policy for a rejected permission, got %#v", names)
	}
}

func TestPolicyManagerErrorsWhenEnforcerIsMissing(t *testing.T) {
	pm := &PolicyManager{}
	roleID := uuid.NewString()
	userID := uuid.NewString()

	if err := pm.AddRoleForUser(userID, roleID, OrgPlatform); err == nil {
		t.Error("AddRoleForUser: expected an error when the enforcer is missing, got nil")
	}
	if err := pm.ReplaceRolesForUser(userID, []string{roleID}, OrgPlatform); err == nil {
		t.Error("ReplaceRolesForUser: expected an error when the enforcer is missing, got nil")
	}
	if err := pm.ReplacePermissionsForRole(roleID, nil, OrgPlatform); err == nil {
		t.Error("ReplacePermissionsForRole: expected an error when the enforcer is missing, got nil")
	}
	if err := pm.AddTranslationGroupOwner(userID, uuid.NewString()); err == nil {
		t.Error("AddTranslationGroupOwner: expected an error when the enforcer is missing, got nil")
	}
	if err := pm.ReplaceBaselinePolicies(); err == nil {
		t.Error("ReplaceBaselinePolicies: expected an error when the enforcer is missing, got nil")
	}
	if _, err := pm.RolesForUser(userID, OrgPlatform); err == nil {
		t.Error("RolesForUser: expected an error when the enforcer is missing, got nil")
	}
	if _, err := pm.PermissionNamesForRole(roleID, OrgPlatform); err == nil {
		t.Error("PermissionNamesForRole: expected an error when the enforcer is missing, got nil")
	}
}

func TestPolicyManagerRejectsEmptyIdentifiers(t *testing.T) {
	pm, _ := newTestPolicyManager(t)

	if err := pm.AddRoleForUser("", uuid.NewString(), OrgPlatform); err == nil {
		t.Error("AddRoleForUser: expected an error for an empty user id, got nil")
	}
	if err := pm.AddRoleForUser(uuid.NewString(), "", OrgPlatform); err == nil {
		t.Error("AddRoleForUser: expected an error for an empty role id, got nil")
	}
	if err := pm.ReplaceRolesForUser(uuid.NewString(), []string{uuid.NewString()}, ""); err == nil {
		t.Error("ReplaceRolesForUser: expected an error for an empty org, got nil")
	}
	if err := pm.ReplacePermissionsForRole("", []string{"comic:read"}, OrgPlatform); err == nil {
		t.Error("ReplacePermissionsForRole: expected an error for an empty role id, got nil")
	}
	if err := pm.AddTranslationGroupMember(uuid.NewString(), ""); err == nil {
		t.Error("AddTranslationGroupMember: expected an error for an empty group id, got nil")
	}
}
