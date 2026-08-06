package authorization

import (
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

func (m *PolicyManager) AddRoleForUser(userID string, roleID string, org Org) error {
	if m == nil || m.enforcer == nil || userID == "" || roleID == "" || org == "" {
		return nil
	}

	_, err := m.enforcer.AddGroupingPolicy(userID, roleID, string(org))
	return err
}

func (m *PolicyManager) RemoveRoleForUser(userID string, roleID string, org Org) error {
	if m == nil || m.enforcer == nil || userID == "" || roleID == "" || org == "" {
		return nil
	}

	_, err := m.enforcer.RemoveGroupingPolicy(userID, roleID, string(org))
	return err
}

func (m *PolicyManager) ReplaceRolesForUser(userID string, roleIDs []string, org Org) error {
	if m == nil || m.enforcer == nil || userID == "" || org == "" {
		return nil
	}

	if _, err := m.enforcer.RemoveFilteredGroupingPolicy(0, userID, "", string(org)); err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		if roleID == "" {
			continue
		}
		if _, err := m.enforcer.AddGroupingPolicy(userID, roleID, string(org)); err != nil {
			return err
		}
	}
	return nil
}

func (m *PolicyManager) RemoveRole(roleID string, org Org) error {
	if m == nil || m.enforcer == nil || roleID == "" || org == "" {
		return nil
	}

	if _, err := m.enforcer.RemoveFilteredGroupingPolicy(0, roleID, "", string(org)); err != nil {
		return err
	}
	_, err := m.enforcer.RemoveFilteredGroupingPolicy(1, roleID, string(org))
	return err
}

func (m *PolicyManager) ReplacePermissionsForRole(roleID string, permissions []PermissionRule, org Org) error {
	if m == nil || m.enforcer == nil || roleID == "" || org == "" {
		return nil
	}

	if _, err := m.enforcer.RemoveFilteredGroupingPolicy(0, roleID, "", string(org)); err != nil {
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
	if m == nil || m.enforcer == nil || roleID == "" || permission.ID == "" || org == "" {
		return nil
	}

	actions, object, ok := parsePermissionName(permission.Name)
	if !ok {
		return nil
	}
	if _, err := m.enforcer.AddGroupingPolicy(roleID, permission.ID, string(org)); err != nil {
		return err
	}
	return m.addPermissionPolicy(org, permission.ID, actions, object, CtxAny)
}

func (m *PolicyManager) RemovePermissionForRole(roleID string, permissionID string, org Org) error {
	if m == nil || m.enforcer == nil || roleID == "" || permissionID == "" || org == "" {
		return nil
	}

	_, err := m.enforcer.RemoveGroupingPolicy(roleID, permissionID, string(org))
	return err
}

func (m *PolicyManager) ReplacePermission(permissionID string, permissionName string, org Org) error {
	if m == nil || m.enforcer == nil || permissionID == "" || org == "" {
		return nil
	}

	actions, object, ok := parsePermissionName(permissionName)
	if !ok {
		return nil
	}
	if _, err := m.enforcer.RemoveFilteredPolicy(0, string(org), permissionID); err != nil {
		return err
	}
	return m.addPermissionPolicy(org, permissionID, actions, object, CtxAny)
}

func (m *PolicyManager) RemovePermission(permissionID string, org Org) error {
	if m == nil || m.enforcer == nil || permissionID == "" || org == "" {
		return nil
	}

	if _, err := m.enforcer.RemoveFilteredGroupingPolicy(1, permissionID, string(org)); err != nil {
		return err
	}
	_, err := m.enforcer.RemoveFilteredPolicy(0, string(org), permissionID)
	return err
}

func (m *PolicyManager) AddTranslationGroupMember(userID string, groupID string) error {
	org := TranslationGroupOrgString(groupID)
	if err := m.ensureTranslationGroupPermissions(org); err != nil {
		return err
	}
	return m.AddRoleForUser(userID, roleGroupMember, org)
}

func (m *PolicyManager) AddTranslationGroupOwner(userID string, groupID string) error {
	org := TranslationGroupOrgString(groupID)
	if err := m.ensureTranslationGroupPermissions(org); err != nil {
		return err
	}
	return m.AddRoleForUser(userID, roleGroupOwner, org)
}

func (m *PolicyManager) RemoveTranslationGroupOwner(userID string, groupID string) error {
	return m.RemoveRoleForUser(userID, roleGroupOwner, TranslationGroupOrgString(groupID))
}

func (m *PolicyManager) ensureTranslationGroupPermissions(org Org) error {
	if m == nil || m.enforcer == nil || org == "" {
		return nil
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
	if _, err := m.enforcer.AddGroupingPolicy(roleID, policy.ID, string(org)); err != nil {
		return err
	}
	return m.addPermissionPolicy(org, policy.ID, policy.Actions, policy.Object, policy.Context)
}

func (m *PolicyManager) addPermissionPolicy(org Org, permissionID string, actions []Action, object Object, ctx Context) error {
	if ctx == "" {
		ctx = CtxAny
	}
	for _, action := range actions {
		if _, err := m.enforcer.AddPolicy(string(org), permissionID, string(action), string(object), string(ctx), "allow"); err != nil {
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

func parsePermissionName(name string) (actions []Action, object Object, ok bool) {
	parts := strings.Split(name, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, "", false
	}

	switch parts[1] {
	case "write":
		return []Action{ActionCreate, ActionUpdate, ActionPublish}, Object(parts[0]), true
	case string(ActionManage):
		return []Action{ActionManage}, Object(parts[0]), true
	default:
		return []Action{Action(parts[1])}, Object(parts[0]), true
	}
}
