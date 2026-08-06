package authorization

import (
	"testing"

	"github.com/google/uuid"
)

func TestPolicyManagerRejectsInvalidPermissionName(t *testing.T) {
	pm, _ := newTestPolicyManager(t)

	err := pm.AddPermissionForRole(uuid.New().String(), PermissionRule{
		ID:   uuid.New().String(),
		Name: "comics.read", // dot instead of colon: not a permission this system understands
	}, OrgPlatform)

	if err == nil {
		t.Fatal("expected an unparsable permission name to be rejected, got silent success")
	}
}

func TestPolicyManagerLeavesNoPolicyBehindForInvalidPermissionName(t *testing.T) {
	pm, _ := newTestPolicyManager(t)
	roleID := uuid.New().String()

	_ = pm.AddPermissionForRole(roleID, PermissionRule{
		ID:   uuid.New().String(),
		Name: "comics.read",
	}, OrgPlatform)

	if roles := pm.enforcer.GetRolesForUserInDomain(roleID, string(OrgPlatform)); len(roles) != 0 {
		t.Fatalf("expected no grouping policy for a rejected permission, got %#v", roles)
	}
}

func TestPolicyManagerRejectsInvalidPermissionNameOnReplace(t *testing.T) {
	pm, _ := newTestPolicyManager(t)

	err := pm.ReplacePermission(uuid.New().String(), "manage-everything", OrgPlatform)
	if err == nil {
		t.Fatal("expected an unparsable permission name to be rejected, got silent success")
	}
}

func TestPolicyManagerErrorsWhenEnforcerIsMissing(t *testing.T) {
	pm := &PolicyManager{}
	roleID := uuid.New().String()
	userID := uuid.New().String()

	if err := pm.AddRoleForUser(userID, roleID, OrgPlatform); err == nil {
		t.Error("AddRoleForUser: expected an error when the enforcer is missing, got nil")
	}
	if err := pm.ReplaceRolesForUser(userID, []string{roleID}, OrgPlatform); err == nil {
		t.Error("ReplaceRolesForUser: expected an error when the enforcer is missing, got nil")
	}
	if err := pm.ReplacePermissionsForRole(roleID, nil, OrgPlatform); err == nil {
		t.Error("ReplacePermissionsForRole: expected an error when the enforcer is missing, got nil")
	}
	if err := pm.AddTranslationGroupOwner(userID, uuid.New().String()); err == nil {
		t.Error("AddTranslationGroupOwner: expected an error when the enforcer is missing, got nil")
	}
}

func TestPolicyManagerRejectsEmptyIdentifiers(t *testing.T) {
	pm, _ := newTestPolicyManager(t)

	if err := pm.AddRoleForUser("", uuid.New().String(), OrgPlatform); err == nil {
		t.Error("AddRoleForUser: expected an error for an empty user id, got nil")
	}
	if err := pm.AddRoleForUser(uuid.New().String(), "", OrgPlatform); err == nil {
		t.Error("AddRoleForUser: expected an error for an empty role id, got nil")
	}
	if err := pm.ReplaceRolesForUser(uuid.New().String(), []string{uuid.New().String()}, ""); err == nil {
		t.Error("ReplaceRolesForUser: expected an error for an empty org, got nil")
	}
}

func TestValidatePermissionNameAcceptsSeededNames(t *testing.T) {
	for _, name := range []string{"comic:read", "comic:write", "role:manage", "translation_group:delete"} {
		if err := ValidatePermissionName(name); err != nil {
			t.Errorf("expected %q to be a valid permission name, got: %v", name, err)
		}
	}
}

func TestValidatePermissionNameRejectsMalformedNames(t *testing.T) {
	for _, name := range []string{"", "comic", "comic:", ":read", "comic:read:extra", "comics.read"} {
		if err := ValidatePermissionName(name); err == nil {
			t.Errorf("expected %q to be rejected as a permission name, got nil", name)
		}
	}
}
