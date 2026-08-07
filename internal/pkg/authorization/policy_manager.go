package authorization

import (
	"errors"
	"fmt"
	"sort"

	casbinpkg "manga-go/internal/pkg/casbin"

	"go.uber.org/fx"
)

const (
	roleGroupOwner  = "group_owner"
	roleGroupMember = "group_member"

	effectAllow = "allow"
	effectDeny  = "deny"
)

// Writing a policy must never fail silently: a caller that believes it granted
// or revoked access, when nothing was written, leaves the database and the
// policy engine permanently out of step.
var (
	ErrPolicyEngineUnavailable = errors.New("policy engine unavailable")
	ErrInvalidPolicyArgument   = errors.New("invalid policy argument")
	ErrInvalidPermissionName   = errors.New("invalid permission name")
)

// PolicyRule is a single row of the policy table: the subject it is attached to
// may perform Action on Object, but only in Context.
type PolicyRule struct {
	Action  Action
	Object  Object
	Context Context
	Deny    bool
}

func (r PolicyRule) effect() string {
	if r.Deny {
		return effectDeny
	}
	return effectAllow
}

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

// ---------------------------------------------------------------- user ↔ role

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

// RolesForUser reports the role ids a user holds. The policy engine is the
// record of who holds what, so this is the only place to ask.
func (m *PolicyManager) RolesForUser(userID string, org Org) ([]string, error) {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return nil, err
	}
	if err := requireNonEmpty("RolesForUser", map[string]string{
		"user id": userID,
		"org":     string(org),
	}); err != nil {
		return nil, err
	}

	return enforcer.GetRolesForUserInDomain(userID, string(org)), nil
}

func (m *PolicyManager) UsersForRole(roleID string, org Org) ([]string, error) {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return nil, err
	}
	if err := requireNonEmpty("UsersForRole", map[string]string{
		"role id": roleID,
		"org":     string(org),
	}); err != nil {
		return nil, err
	}

	users := enforcer.GetUsersForRoleInDomain(roleID, string(org))
	sort.Strings(users)
	return users, nil
}

// RemoveRole deletes a role: both what it grants and who holds it.
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

	if _, err := enforcer.RemoveFilteredPolicy(0, string(org), roleID); err != nil {
		return err
	}
	_, err = enforcer.RemoveFilteredGroupingPolicy(1, roleID, string(org))
	return err
}

// ---------------------------------------------------------- role ↔ permission

// ReplacePermissionsForRole sets the role's grants to exactly names. Every name
// is checked against the catalog first: a batch containing one unknown name
// changes nothing, so a typo cannot strip a role of everything it had.
func (m *PolicyManager) ReplacePermissionsForRole(roleID string, names []string, org Org) error {
	if _, err := m.requireEnforcer(); err != nil {
		return err
	}
	if err := requireNonEmpty("ReplacePermissionsForRole", map[string]string{
		"role id": roleID,
		"org":     string(org),
	}); err != nil {
		return err
	}

	rules, err := rulesForPermissionNames(names)
	if err != nil {
		return err
	}
	return m.replacePoliciesForSubject(roleID, org, rules)
}

func (m *PolicyManager) GrantPermissionToRole(roleID string, name string, org Org) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}
	if err := requireNonEmpty("GrantPermissionToRole", map[string]string{
		"role id": roleID,
		"org":     string(org),
	}); err != nil {
		return err
	}

	rules, err := rulesForPermissionNames([]string{name})
	if err != nil {
		return err
	}
	return addPolicies(enforcer, roleID, org, rules)
}

func (m *PolicyManager) RevokePermissionFromRole(roleID string, name string, org Org) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}
	if err := requireNonEmpty("RevokePermissionFromRole", map[string]string{
		"role id": roleID,
		"org":     string(org),
	}); err != nil {
		return err
	}

	rules, err := rulesForPermissionNames([]string{name})
	if err != nil {
		return err
	}
	for _, rule := range rules {
		if _, err := enforcer.RemoveFilteredPolicy(
			0, string(org), roleID, string(rule.Action), string(rule.Object), string(rule.Context),
		); err != nil {
			return err
		}
	}
	return nil
}

// PermissionNamesForRole reads the role's grants back as catalog names. The
// expansion of a shorthand is folded back up, so a role granted "comic:write"
// reads back as "comic:write" rather than three separate actions.
func (m *PolicyManager) PermissionNamesForRole(roleID string, org Org) ([]string, error) {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return nil, err
	}
	if err := requireNonEmpty("PermissionNamesForRole", map[string]string{
		"role id": roleID,
		"org":     string(org),
	}); err != nil {
		return nil, err
	}

	rows, err := enforcer.GetFilteredPolicy(0, string(org), roleID)
	if err != nil {
		return nil, err
	}

	granted := map[Object]map[Action]bool{}
	for _, row := range rows {
		if len(row) < 6 || row[4] != string(CtxAny) || row[5] != effectAllow {
			continue
		}
		object, action := Object(row[3]), Action(row[2])
		if granted[object] == nil {
			granted[object] = map[Action]bool{}
		}
		granted[object][action] = true
	}

	names := make([]string, 0, len(granted))
	for object, actions := range granted {
		names = append(names, foldActionsToNames(object, actions)...)
	}
	sort.Strings(names)
	return names, nil
}

// ------------------------------------------------------------------ baseline

// ReplaceBaselinePolicies writes the rights every visitor has before any role
// is considered: what an anonymous visitor may read, and what merely being
// signed in allows. Keeping these as policy rows rather than Go branches means
// enforcement has exactly one source of truth.
func (m *PolicyManager) ReplaceBaselinePolicies() error {
	if _, err := m.requireEnforcer(); err != nil {
		return err
	}

	for _, subject := range []string{SubjectAnonymous, RoleAuthenticated} {
		if err := m.replacePoliciesForSubject(subject, OrgPlatform, baselinePolicies()[subject]); err != nil {
			return err
		}
	}
	return nil
}

func baselinePolicies() map[string][]PolicyRule {
	return map[string][]PolicyRule{
		SubjectAnonymous: {
			{Action: ActionRead, Object: ObjectComic, Context: CtxPublished},
			{Action: ActionRead, Object: ObjectChapter, Context: CtxPublished},
		},
		RoleAuthenticated: {
			{Action: ActionCreate, Object: ObjectComment, Context: CtxAny},
			{Action: ActionUpdate, Object: ObjectComment, Context: CtxOwner},
			{Action: ActionDelete, Object: ObjectComment, Context: CtxOwner},

			{Action: ActionCreate, Object: ObjectRating, Context: CtxAny},
			{Action: ActionUpdate, Object: ObjectRating, Context: CtxOwner},
			{Action: ActionDelete, Object: ObjectRating, Context: CtxOwner},

			{Action: ActionCreate, Object: ObjectReadingHistory, Context: CtxAny},
			{Action: ActionRead, Object: ObjectReadingHistory, Context: CtxOwner},
			{Action: ActionUpdate, Object: ObjectReadingHistory, Context: CtxOwner},
			{Action: ActionDelete, Object: ObjectReadingHistory, Context: CtxOwner},

			{Action: ActionUpdate, Object: ObjectUser, Context: CtxSelf},
		},
	}
}

// ---------------------------------------------------------- translation group

func (m *PolicyManager) AddTranslationGroupMember(userID string, groupID string) error {
	if err := requireNonEmpty("AddTranslationGroupMember", map[string]string{"group id": groupID}); err != nil {
		return err
	}

	org := TranslationGroupOrgString(groupID)
	if err := m.ensureTranslationGroupPolicies(org); err != nil {
		return err
	}
	return m.AddRoleForUser(userID, roleGroupMember, org)
}

func (m *PolicyManager) AddTranslationGroupOwner(userID string, groupID string) error {
	if err := requireNonEmpty("AddTranslationGroupOwner", map[string]string{"group id": groupID}); err != nil {
		return err
	}

	org := TranslationGroupOrgString(groupID)
	if err := m.ensureTranslationGroupPolicies(org); err != nil {
		return err
	}
	return m.AddRoleForUser(userID, roleGroupOwner, org)
}

func (m *PolicyManager) RemoveTranslationGroupMember(userID string, groupID string) error {
	if err := requireNonEmpty("RemoveTranslationGroupMember", map[string]string{"group id": groupID}); err != nil {
		return err
	}
	return m.RemoveRoleForUser(userID, roleGroupMember, TranslationGroupOrgString(groupID))
}

func (m *PolicyManager) RemoveTranslationGroupOwner(userID string, groupID string) error {
	if err := requireNonEmpty("RemoveTranslationGroupOwner", map[string]string{"group id": groupID}); err != nil {
		return err
	}
	return m.RemoveRoleForUser(userID, roleGroupOwner, TranslationGroupOrgString(groupID))
}

func (m *PolicyManager) ensureTranslationGroupPolicies(org Org) error {
	if _, err := m.requireEnforcer(); err != nil {
		return err
	}
	if err := requireNonEmpty("ensureTranslationGroupPolicies", map[string]string{"org": string(org)}); err != nil {
		return err
	}

	if err := m.replacePoliciesForSubject(roleGroupMember, org, translationGroupMemberPolicies()); err != nil {
		return err
	}
	return m.replacePoliciesForSubject(roleGroupOwner, org, translationGroupOwnerPolicies())
}

func translationGroupMemberPolicies() []PolicyRule {
	return []PolicyRule{
		{Action: ActionCreate, Object: ObjectChapter, Context: CtxGroupMember},
		{Action: ActionUpdate, Object: ObjectChapter, Context: CtxOwner},
		{Action: ActionPublish, Object: ObjectChapter, Context: CtxGroupMember},
	}
}

func translationGroupOwnerPolicies() []PolicyRule {
	return []PolicyRule{
		{Action: ActionManage, Object: ObjectTranslationGroup, Context: CtxGroupOwner},
		{Action: ActionManage, Object: ObjectComic, Context: CtxGroupMember},
		{Action: ActionManage, Object: ObjectChapter, Context: CtxGroupMember},
	}
}

// -------------------------------------------------------------------- shared

func (m *PolicyManager) replacePoliciesForSubject(subject string, org Org, rules []PolicyRule) error {
	enforcer, err := m.requireEnforcer()
	if err != nil {
		return err
	}

	if _, err := enforcer.RemoveFilteredPolicy(0, string(org), subject); err != nil {
		return err
	}
	return addPolicies(enforcer, subject, org, rules)
}

func addPolicies(enforcer *casbinpkg.Enforcer, subject string, org Org, rules []PolicyRule) error {
	for _, rule := range rules {
		ctx := rule.Context
		if ctx == "" {
			ctx = CtxAny
		}
		if _, err := enforcer.AddPolicy(
			string(org), subject, string(rule.Action), string(rule.Object), string(ctx), rule.effect(),
		); err != nil {
			return err
		}
	}
	return nil
}

func rulesForPermissionNames(names []string) ([]PolicyRule, error) {
	rules := make([]PolicyRule, 0, len(names))
	for _, name := range names {
		def, ok := LookupPermission(name)
		if !ok {
			return nil, fmt.Errorf("%w: %q is not a known permission", ErrInvalidPermissionName, name)
		}
		for _, action := range def.Grants {
			rules = append(rules, PolicyRule{Action: action, Object: def.Object, Context: CtxAny})
		}
	}
	return rules, nil
}

// foldActionsToNames is the inverse of the catalog expansion: it prefers the
// broadest name that the granted actions fully cover.
func foldActionsToNames(object Object, actions map[Action]bool) []string {
	names := []string{}

	for _, candidate := range grantableActions[object] {
		expansion := expandAction(candidate)
		if len(expansion) < 2 {
			continue
		}
		if !coversAll(actions, expansion) {
			continue
		}
		names = append(names, fmt.Sprintf("%s:%s", object, candidate))
		for _, action := range expansion {
			delete(actions, action)
		}
	}

	for action := range actions {
		names = append(names, fmt.Sprintf("%s:%s", object, action))
	}
	return names
}

func coversAll(actions map[Action]bool, required []Action) bool {
	for _, action := range required {
		if !actions[action] {
			return false
		}
	}
	return true
}
