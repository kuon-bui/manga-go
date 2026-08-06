package authorization

import (
	"errors"
	"fmt"
	"strings"

	casbinpkg "manga-go/internal/pkg/casbin"

	"go.uber.org/fx"
)

const (
	roleGroupOwner  = "group_owner"
	roleGroupMember = "group_member"

	permissionGroupMemberCreateChapter  = "group_member:chapter:create"
	permissionGroupMemberUpdateChapter  = "group_member:chapter:update"
	permissionGroupMemberPublishChapter = "group_member:chapter:publish"

	permissionGroupOwnerManageTranslationGroup = "group_owner:translation_group:manage"
	permissionGroupOwnerManageComic            = "group_owner:comic:manage"
	permissionGroupOwnerManageChapter          = "group_owner:chapter:manage"
)

// Writing a policy must never fail silently: a caller that believes it granted
// or revoked access, when nothing was written, leaves the database and the
// policy engine permanently out of step.
var (
	ErrPolicyEngineUnavailable = errors.New("policy engine unavailable")
	ErrInvalidPolicyArgument   = errors.New("invalid policy argument")
	ErrInvalidPermissionName   = errors.New("invalid permission name")
)

type PolicyManager struct {
	enforcer *casbinpkg.Enforcer
}

type PermissionRule struct {
	ID   string
	Name string
}

type permissionPolicy struct {
	ID      string
	Actions []Action
	Object  Object
	Context Context
}

type PolicyManagerParams struct {
	fx.In

	Enforcer *casbinpkg.Enforcer
}

func NewPolicyManager(p PolicyManagerParams) *PolicyManager {
	return &PolicyManager{enforcer: p.Enforcer}
}

func (m *PolicyManager) requireEnforcer() (*casbinpkg.Enforcer, error) {
	if m == nil || m.enforcer == nil {
		return nil, ErrPolicyEngineUnavailable
	}
	return m.enforcer, nil
}

func requireNonEmpty(operation string, fields map[string]string) error {
	for name, value := range fields {
		if value == "" {
			return fmt.Errorf("%w: %s requires a non-empty %s", ErrInvalidPolicyArgument, operation, name)
		}
	}
	return nil
}

func (m *PolicyManager) AddRoleForUser(userID string, roleID string, org Org) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}
	if err := requireNonEmpty("AddRoleForUser", map[string]string{
		"user id": userID,
		"role id": roleID,
		"org":     string(org),
	}); err != nil {
		return err
	}

	_, err = enforcer.AddGroupingPolicy(userID, roleID, string(org))
	return err
}

func (m *PolicyManager) RemoveRoleForUser(userID string, roleID string, org Org) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}
	if err := requireNonEmpty("RemoveRoleForUser", map[string]string{
		"user id": userID,
		"role id": roleID,
		"org":     string(org),
	}); err != nil {
		return err
	}

	_, err = enforcer.RemoveGroupingPolicy(userID, roleID, string(org))
	return err
}

func (m *PolicyManager) ReplaceRolesForUser(userID string, roleIDs []string, org Org) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}
	if err := requireNonEmpty("ReplaceRolesForUser", map[string]string{
		"user id": userID,
		"org":     string(org),
	}); err != nil {
		return err
	}

	if _, err := enforcer.RemoveFilteredGroupingPolicy(0, userID, "", string(org)); err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		if err := m.AddRoleForUser(userID, roleID, org); err != nil {
			return err
		}
	}
	return nil
}

func (m *PolicyManager) RemoveRole(roleID string, org Org) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}
	if err := requireNonEmpty("RemoveRole", map[string]string{
		"role id": roleID,
		"org":     string(org),
	}); err != nil {
		return err
	}

	if _, err := enforcer.RemoveFilteredGroupingPolicy(0, roleID, "", string(org)); err != nil {
		return err
	}
	_, err = enforcer.RemoveFilteredGroupingPolicy(1, roleID, string(org))
	return err
}

func (m *PolicyManager) ReplacePermissionsForRole(roleID string, permissions []PermissionRule, org Org) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}
	if err := requireNonEmpty("ReplacePermissionsForRole", map[string]string{
		"role id": roleID,
		"org":     string(org),
	}); err != nil {
		return err
	}

	// Reject the whole batch before touching the engine, so one bad name cannot
	// leave the role stripped of every permission it used to have.
	for _, permission := range permissions {
		if err := ValidatePermissionName(permission.Name); err != nil {
			return err
		}
	}

	if _, err := enforcer.RemoveFilteredGroupingPolicy(0, roleID, "", string(org)); err != nil {
		return err
	}
	for _, permission := range permissions {
		if err := m.AddPermissionForRole(roleID, permission, org); err != nil {
			return err
		}
	}
	return nil
}

func (m *PolicyManager) AddPermissionForRole(roleID string, permission PermissionRule, org Org) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}
	if err := requireNonEmpty("AddPermissionForRole", map[string]string{
		"role id":       roleID,
		"permission id": permission.ID,
		"org":           string(org),
	}); err != nil {
		return err
	}

	actions, object, err := parsePermissionName(permission.Name)
	if err != nil {
		return err
	}
	if _, err := enforcer.AddGroupingPolicy(roleID, permission.ID, string(org)); err != nil {
		return err
	}
	return m.addPermissionPolicy(org, permission.ID, actions, object, CtxAny)
}

func (m *PolicyManager) RemovePermissionForRole(roleID string, permissionID string, org Org) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}
	if err := requireNonEmpty("RemovePermissionForRole", map[string]string{
		"role id":       roleID,
		"permission id": permissionID,
		"org":           string(org),
	}); err != nil {
		return err
	}

	_, err = enforcer.RemoveGroupingPolicy(roleID, permissionID, string(org))
	return err
}

func (m *PolicyManager) ReplacePermission(permissionID string, permissionName string, org Org) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}
	if err := requireNonEmpty("ReplacePermission", map[string]string{
		"permission id": permissionID,
		"org":           string(org),
	}); err != nil {
		return err
	}

	actions, object, err := parsePermissionName(permissionName)
	if err != nil {
		return err
	}
	if _, err := enforcer.RemoveFilteredPolicy(0, string(org), permissionID); err != nil {
		return err
	}
	return m.addPermissionPolicy(org, permissionID, actions, object, CtxAny)
}

func (m *PolicyManager) RemovePermission(permissionID string, org Org) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}
	if err := requireNonEmpty("RemovePermission", map[string]string{
		"permission id": permissionID,
		"org":           string(org),
	}); err != nil {
		return err
	}

	if _, err := enforcer.RemoveFilteredGroupingPolicy(1, permissionID, string(org)); err != nil {
		return err
	}
	_, err = enforcer.RemoveFilteredPolicy(0, string(org), permissionID)
	return err
}

func (m *PolicyManager) AddTranslationGroupMember(userID string, groupID string) error {
	if err := requireNonEmpty("AddTranslationGroupMember", map[string]string{
		"group id": groupID,
	}); err != nil {
		return err
	}

	org := TranslationGroupOrgString(groupID)
	if err := m.ensureTranslationGroupPermissions(org); err != nil {
		return err
	}
	return m.AddRoleForUser(userID, roleGroupMember, org)
}

func (m *PolicyManager) AddTranslationGroupOwner(userID string, groupID string) error {
	if err := requireNonEmpty("AddTranslationGroupOwner", map[string]string{
		"group id": groupID,
	}); err != nil {
		return err
	}

	org := TranslationGroupOrgString(groupID)
	if err := m.ensureTranslationGroupPermissions(org); err != nil {
		return err
	}
	return m.AddRoleForUser(userID, roleGroupOwner, org)
}

func (m *PolicyManager) RemoveTranslationGroupOwner(userID string, groupID string) error {
	if err := requireNonEmpty("RemoveTranslationGroupOwner", map[string]string{
		"group id": groupID,
	}); err != nil {
		return err
	}

	return m.RemoveRoleForUser(userID, roleGroupOwner, TranslationGroupOrgString(groupID))
}

func (m *PolicyManager) ensureTranslationGroupPermissions(org Org) error {
	if _, err := m.requireEnforcer(); err != nil {
		return err
	}
	if err := requireNonEmpty("ensureTranslationGroupPermissions", map[string]string{
		"org": string(org),
	}); err != nil {
		return err
	}

	for _, policy := range translationGroupMemberPolicies() {
		if err := m.addPermissionForBuiltinRole(org, roleGroupMember, policy); err != nil {
			return err
		}
	}
	for _, policy := range translationGroupOwnerPolicies() {
		if err := m.addPermissionForBuiltinRole(org, roleGroupOwner, policy); err != nil {
			return err
		}
	}
	return nil
}

func (m *PolicyManager) addPermissionForBuiltinRole(org Org, roleID string, policy permissionPolicy) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}

	if _, err := enforcer.AddGroupingPolicy(roleID, policy.ID, string(org)); err != nil {
		return err
	}
	return m.addPermissionPolicy(org, policy.ID, policy.Actions, policy.Object, policy.Context)
}

func (m *PolicyManager) addPermissionPolicy(org Org, permissionID string, actions []Action, object Object, ctx Context) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}

	if ctx == "" {
		ctx = CtxAny
	}
	for _, action := range actions {
		if _, err := enforcer.AddPolicy(string(org), permissionID, string(action), string(object), string(ctx), "allow"); err != nil {
			return err
		}
	}
	return nil
}

func translationGroupMemberPolicies() []permissionPolicy {
	return []permissionPolicy{
		{
			ID:      permissionGroupMemberCreateChapter,
			Actions: []Action{ActionCreate},
			Object:  ObjectChapter,
			Context: CtxGroupMember,
		},
		{
			ID:      permissionGroupMemberUpdateChapter,
			Actions: []Action{ActionUpdate},
			Object:  ObjectChapter,
			Context: CtxOwner,
		},
		{
			ID:      permissionGroupMemberPublishChapter,
			Actions: []Action{ActionPublish},
			Object:  ObjectChapter,
			Context: CtxGroupMember,
		},
	}
}

func translationGroupOwnerPolicies() []permissionPolicy {
	return []permissionPolicy{
		{
			ID:      permissionGroupOwnerManageTranslationGroup,
			Actions: []Action{ActionManage},
			Object:  ObjectTranslationGroup,
			Context: CtxGroupOwner,
		},
		{
			ID:      permissionGroupOwnerManageComic,
			Actions: []Action{ActionManage},
			Object:  ObjectComic,
			Context: CtxGroupMember,
		},
		{
			ID:      permissionGroupOwnerManageChapter,
			Actions: []Action{ActionManage},
			Object:  ObjectChapter,
			Context: CtxGroupMember,
		},
	}
}

// ValidatePermissionName reports whether name is a permission this system can
// translate into policy rules. Permission names are "<object>:<action>", where
// "write" is shorthand for create + update + publish.
func ValidatePermissionName(name string) error {
	_, _, err := parsePermissionName(name)
	return err
}

func parsePermissionName(name string) ([]Action, Object, error) {
	parts := strings.Split(name, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, "", fmt.Errorf("%w: %q is not in <object>:<action> form", ErrInvalidPermissionName, name)
	}

	switch parts[1] {
	case "write":
		return []Action{ActionCreate, ActionUpdate, ActionPublish}, Object(parts[0]), nil
	case string(ActionManage):
		return []Action{ActionManage}, Object(parts[0]), nil
	default:
		return []Action{Action(parts[1])}, Object(parts[0]), nil
	}
}
