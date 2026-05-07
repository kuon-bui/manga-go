package authorization

import (
	"strings"

	casbinpkg "manga-go/internal/pkg/casbin"

	"go.uber.org/fx"
)

const (
	roleGroupOwner  = "group_owner"
	roleGroupMember = "group_member"
)

type PolicyManager struct {
	enforcer *casbinpkg.Enforcer
}

type PolicyManagerParams struct {
	fx.In

	Enforcer *casbinpkg.Enforcer
}

func NewPolicyManager(p PolicyManagerParams) *PolicyManager {
	return &PolicyManager{enforcer: p.Enforcer}
}

func (m *PolicyManager) AddRoleForUser(userID string, role string, org string) error {
	if m == nil || m.enforcer == nil || userID == "" || role == "" || org == "" {
		return nil
	}

	_, err := m.enforcer.AddGroupingPolicy(userID, role, org)
	return err
}

func (m *PolicyManager) RemoveRoleForUser(userID string, role string, org string) error {
	if m == nil || m.enforcer == nil || userID == "" || role == "" || org == "" {
		return nil
	}

	_, err := m.enforcer.RemoveGroupingPolicy(userID, role, org)
	return err
}

func (m *PolicyManager) ReplaceRolesForUser(userID string, roles []string, org string) error {
	if m == nil || m.enforcer == nil || userID == "" || org == "" {
		return nil
	}

	if _, err := m.enforcer.RemoveFilteredGroupingPolicy(0, userID, "", org); err != nil {
		return err
	}
	for _, role := range roles {
		if role == "" {
			continue
		}
		if _, err := m.enforcer.AddGroupingPolicy(userID, role, org); err != nil {
			return err
		}
	}
	return nil
}

func (m *PolicyManager) RemoveRole(role string, org string) error {
	if m == nil || m.enforcer == nil || role == "" || org == "" {
		return nil
	}

	if _, err := m.enforcer.RemoveFilteredGroupingPolicy(1, role, org); err != nil {
		return err
	}
	_, err := m.enforcer.RemoveFilteredPolicy(0, org, role)
	return err
}

func (m *PolicyManager) RenameRole(oldRole string, newRole string, org string) error {
	if m == nil || m.enforcer == nil || oldRole == "" || newRole == "" || org == "" {
		return nil
	}

	groupingPolicies, err := m.enforcer.GetFilteredGroupingPolicy(1, oldRole, org)
	if err != nil {
		return err
	}
	rolePolicies, err := m.enforcer.GetFilteredPolicy(0, org, oldRole)
	if err != nil {
		return err
	}

	if err := m.RemoveRole(oldRole, org); err != nil {
		return err
	}

	for _, policy := range groupingPolicies {
		if len(policy) < 3 {
			continue
		}
		if _, err := m.enforcer.AddGroupingPolicy(policy[0], newRole, policy[2]); err != nil {
			return err
		}
	}
	for _, policy := range rolePolicies {
		if len(policy) < 6 {
			continue
		}
		policy[1] = newRole
		if _, err := m.enforcer.AddPolicy(policy); err != nil {
			return err
		}
	}
	return nil
}

func (m *PolicyManager) ReplacePermissionsForRole(role string, permissionNames []string, org string) error {
	if m == nil || m.enforcer == nil || role == "" || org == "" {
		return nil
	}

	policies, err := m.enforcer.GetFilteredPolicy(0, org, role)
	if err != nil {
		return err
	}
	for _, policy := range policies {
		if len(policy) < 6 || policy[4] != CtxAny || policy[5] != "allow" {
			continue
		}
		if _, err := m.enforcer.RemovePolicy(policy); err != nil {
			return err
		}
	}
	for _, name := range permissionNames {
		if err := m.AddPermissionForRole(role, name, org); err != nil {
			return err
		}
	}
	return nil
}

func (m *PolicyManager) AddPermissionForRole(role string, permissionName string, org string) error {
	if m == nil || m.enforcer == nil || role == "" || org == "" {
		return nil
	}

	actions, object, ok := parsePermissionName(permissionName)
	if !ok {
		return nil
	}
	for _, action := range actions {
		if _, err := m.enforcer.AddPolicy(org, role, action, object, CtxAny, "allow"); err != nil {
			return err
		}
	}
	return nil
}

func (m *PolicyManager) RemovePermissionForRole(role string, permissionName string, org string) error {
	if m == nil || m.enforcer == nil || role == "" || org == "" {
		return nil
	}

	actions, object, ok := parsePermissionName(permissionName)
	if !ok {
		return nil
	}
	for _, action := range actions {
		if _, err := m.enforcer.RemovePolicy(org, role, action, object, CtxAny, "allow"); err != nil {
			return err
		}
	}
	return nil
}

func (m *PolicyManager) ReplacePermissionName(oldName string, newName string, org string) error {
	if m == nil || m.enforcer == nil || org == "" {
		return nil
	}

	oldActions, oldObject, ok := parsePermissionName(oldName)
	if !ok {
		return nil
	}
	newActions, newObject, ok := parsePermissionName(newName)
	if !ok {
		return nil
	}

	policies, err := m.enforcer.GetFilteredPolicy(0, org)
	if err != nil {
		return err
	}

	affectedRoles := make(map[string]struct{})
	for _, policy := range policies {
		if len(policy) < 6 || policy[3] != oldObject || policy[4] != CtxAny || policy[5] != "allow" {
			continue
		}
		if !contains(oldActions, policy[2]) {
			continue
		}
		if _, err := m.enforcer.RemovePolicy(policy); err != nil {
			return err
		}
		affectedRoles[policy[1]] = struct{}{}
	}

	for role := range affectedRoles {
		for _, action := range newActions {
			if _, err := m.enforcer.AddPolicy(org, role, action, newObject, CtxAny, "allow"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *PolicyManager) RemovePermissionName(permissionName string, org string) error {
	if m == nil || m.enforcer == nil || org == "" {
		return nil
	}

	actions, object, ok := parsePermissionName(permissionName)
	if !ok {
		return nil
	}

	policies, err := m.enforcer.GetFilteredPolicy(0, org)
	if err != nil {
		return err
	}
	for _, policy := range policies {
		if len(policy) < 6 || policy[3] != object || policy[4] != CtxAny || policy[5] != "allow" {
			continue
		}
		if !contains(actions, policy[2]) {
			continue
		}
		if _, err := m.enforcer.RemovePolicy(policy); err != nil {
			return err
		}
	}
	return nil
}

func (m *PolicyManager) AddTranslationGroupMember(userID string, groupID string) error {
	return m.AddRoleForUser(userID, roleGroupMember, "tg:"+groupID)
}

func (m *PolicyManager) AddTranslationGroupOwner(userID string, groupID string) error {
	return m.AddRoleForUser(userID, roleGroupOwner, "tg:"+groupID)
}

func (m *PolicyManager) RemoveTranslationGroupOwner(userID string, groupID string) error {
	return m.RemoveRoleForUser(userID, roleGroupOwner, "tg:"+groupID)
}

func parsePermissionName(name string) (actions []string, object string, ok bool) {
	parts := strings.Split(name, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, "", false
	}

	switch parts[1] {
	case "write":
		return []string{ActionCreate, ActionUpdate, ActionPublish}, parts[0], true
	case ActionManage:
		return []string{ActionManage}, parts[0], true
	default:
		return []string{parts[1]}, parts[0], true
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
