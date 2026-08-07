package authorizationadmin

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	codeAuthorizationStateChanged = "AUTHORIZATION_STATE_CHANGED"
	codeSelfManageRequired        = "SELF_MANAGE_REQUIRED"
	codeLastRoleManager           = "LAST_ROLE_MANAGER"
)

type roleAssignmentSnapshot struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Service) ReplaceUserRoles(
	ctx context.Context,
	userID uuid.UUID,
	roleIDs []uuid.UUID,
	expectedVersion string,
) response.Result {
	return s.locker.WithLock(ctx, func() response.Result {
		return s.replaceUserRolesLocked(ctx, userID, roleIDs, expectedVersion)
	})
}

func (s *Service) replaceUserRolesLocked(
	ctx context.Context,
	userID uuid.UUID,
	roleIDs []uuid.UUID,
	expectedVersion string,
) response.Result {
	target, err := s.userRepo.FindOne(ctx, []any{clause.Eq{Column: "id", Value: userID}}, nil)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.ResultNotFound("User")
	}
	if err != nil {
		return s.dbFailure("find authorization target user", err)
	}

	roleIDs = uniqueUUIDs(roleIDs)
	for _, roleID := range roleIDs {
		if _, err := s.roleRepo.FindOne(ctx, []any{clause.Eq{Column: "id", Value: roleID}}, nil); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.ResultNotFound("Role")
			}
			return s.dbFailure("validate authorization role", err)
		}
	}

	globalVersion, userVersion, err := s.revisions.Current(ctx, userID)
	if err != nil {
		return s.dbFailure("read authorization revision", err)
	}
	currentVersion := fmt.Sprintf("g%d:u%d", globalVersion, userVersion)
	if expectedVersion != "" && expectedVersion != currentVersion {
		return authorizationConflict(codeAuthorizationStateChanged, "Authorization state changed; refresh and try again")
	}

	before, err := s.policyManager.RolesForUser(userID.String(), authorization.OrgPlatform)
	if err != nil {
		return s.internalFailure("read current user roles", err)
	}
	sort.Strings(before)
	beforeSnapshot, err := s.roleAssignmentSnapshots(ctx, before)
	if err != nil {
		return s.dbFailure("snapshot current user roles", err)
	}
	beforeManagers, err := s.roleManagerCount(ctx)
	if err != nil {
		return s.internalFailure("count role managers", err)
	}
	viewer := authorization.ViewerFromContext(ctx)
	actorHadManage, err := s.subjectCanManageRoles(ctx, viewer.Subject)
	if err != nil {
		return s.internalFailure("check actor authorization", err)
	}

	after := stringifyUUIDs(roleIDs)
	afterSnapshot, err := s.roleAssignmentSnapshots(ctx, after)
	if err != nil {
		return s.dbFailure("snapshot replacement user roles", err)
	}
	if err := s.policyManager.ReplaceRolesForUser(userID.String(), after, authorization.OrgPlatform); err != nil {
		_ = s.policyManager.ReplaceRolesForUser(userID.String(), before, authorization.OrgPlatform)
		return s.internalFailure("replace user roles", err)
	}
	rollback := func() {
		if err := s.policyManager.ReplaceRolesForUser(userID.String(), before, authorization.OrgPlatform); err != nil {
			s.logMutationError("restore user roles", err)
		}
	}

	if viewer.User != nil && viewer.User.ID == userID && actorHadManage {
		actorHasManage, err := s.subjectCanManageRoles(ctx, viewer.Subject)
		if err != nil {
			rollback()
			return s.internalFailure("check actor authorization after mutation", err)
		}
		if !actorHasManage {
			rollback()
			return authorizationConflict(codeSelfManageRequired, "You cannot remove your own role-management access")
		}
	}
	afterManagers, err := s.roleManagerCount(ctx)
	if err != nil {
		rollback()
		return s.internalFailure("count role managers after mutation", err)
	}
	if beforeManagers > 0 && afterManagers == 0 {
		rollback()
		return authorizationConflict(codeLastRoleManager, "At least one user must retain role-management access")
	}

	entry := newAuthorizationAudit(viewer, "user.roles_replaced", "user", target.ID, target.Name,
		common.JSONMap{"roles": beforeSnapshot}, common.JSONMap{"roles": afterSnapshot})
	if result := s.persistMutation(entry, func(tx *gorm.DB) error {
		_, err := s.revisions.BumpUserTx(tx, userID)
		return err
	}); result.IsError() {
		rollback()
		return result
	}

	if s.cache != nil {
		key := profileCacheKey(userID, globalVersion, userVersion)
		if err := s.cache.Delete(ctx, key); err != nil {
			s.logCacheError("delete", key, err)
		}
	}
	return response.ResultSuccess("Roles assigned successfully", map[string]any{
		"roleIds": after,
		"version": fmt.Sprintf("g%d:u%d", globalVersion, userVersion+1),
	})
}

func (s *Service) roleAssignmentSnapshots(
	ctx context.Context,
	roleIDs []string,
) ([]roleAssignmentSnapshot, error) {
	snapshots := make([]roleAssignmentSnapshot, 0, len(roleIDs))
	for _, rawID := range roleIDs {
		roleID, err := uuid.Parse(rawID)
		if err != nil {
			return nil, fmt.Errorf("parse role snapshot ID %q: %w", rawID, err)
		}
		role, err := s.roleRepo.FindOne(ctx, []any{clause.Eq{Column: "id", Value: roleID}}, nil)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, roleAssignmentSnapshot{ID: rawID, Name: role.Name})
	}
	sort.Slice(snapshots, func(i, j int) bool {
		if snapshots[i].Name == snapshots[j].Name {
			return snapshots[i].ID < snapshots[j].ID
		}
		return snapshots[i].Name < snapshots[j].Name
	})
	return snapshots, nil
}

func (s *Service) ReplaceRolePermissions(
	ctx context.Context,
	roleID uuid.UUID,
	permissions []string,
	expectedVersion string,
) response.Result {
	return s.locker.WithLock(ctx, func() response.Result {
		return s.replaceRolePermissionsLocked(ctx, roleID, permissions, expectedVersion)
	})
}

func (s *Service) replaceRolePermissionsLocked(
	ctx context.Context,
	roleID uuid.UUID,
	permissions []string,
	expectedVersion string,
) response.Result {
	role, err := s.roleRepo.FindOne(ctx, []any{clause.Eq{Column: "id", Value: roleID}}, nil)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.ResultNotFound("Role")
	}
	if err != nil {
		return s.dbFailure("find authorization role", err)
	}

	permissions = uniqueStrings(permissions)
	for _, name := range permissions {
		if err := authorization.ValidatePermissionName(name); err != nil {
			return response.ResultError(err.Error())
		}
	}
	globalVersion, err := s.revisions.CurrentGlobal(ctx)
	if err != nil {
		return s.dbFailure("read authorization revision", err)
	}
	currentVersion := fmt.Sprintf("g%d", globalVersion)
	if expectedVersion != "" && expectedVersion != currentVersion {
		return authorizationConflict(codeAuthorizationStateChanged, "Authorization state changed; refresh and try again")
	}

	before, err := s.policyManager.PermissionNamesForRole(roleID.String(), authorization.OrgPlatform)
	if err != nil {
		return s.internalFailure("read current role permissions", err)
	}
	sort.Strings(before)
	beforeManagers, err := s.roleManagerCount(ctx)
	if err != nil {
		return s.internalFailure("count role managers", err)
	}
	viewer := authorization.ViewerFromContext(ctx)
	actorHadManage, err := s.subjectCanManageRoles(ctx, viewer.Subject)
	if err != nil {
		return s.internalFailure("check actor authorization", err)
	}

	if err := s.policyManager.ReplacePermissionsForRole(roleID.String(), permissions, authorization.OrgPlatform); err != nil {
		_ = s.policyManager.ReplacePermissionsForRole(roleID.String(), before, authorization.OrgPlatform)
		return s.internalFailure("replace role permissions", err)
	}
	rollback := func() {
		if err := s.policyManager.ReplacePermissionsForRole(roleID.String(), before, authorization.OrgPlatform); err != nil {
			s.logMutationError("restore role permissions", err)
		}
	}

	if actorHadManage {
		actorHasManage, err := s.subjectCanManageRoles(ctx, viewer.Subject)
		if err != nil {
			rollback()
			return s.internalFailure("check actor authorization after mutation", err)
		}
		if !actorHasManage {
			rollback()
			return authorizationConflict(codeSelfManageRequired, "You cannot remove your own role-management access")
		}
	}
	afterManagers, err := s.roleManagerCount(ctx)
	if err != nil {
		rollback()
		return s.internalFailure("count role managers after mutation", err)
	}
	if beforeManagers > 0 && afterManagers == 0 {
		rollback()
		return authorizationConflict(codeLastRoleManager, "At least one user must retain role-management access")
	}

	entry := newAuthorizationAudit(viewer, "role.permissions_replaced", "role", role.ID, role.Name,
		common.JSONMap{"permissions": before}, common.JSONMap{"permissions": permissions})
	if result := s.persistMutation(entry, func(tx *gorm.DB) error {
		_, err := s.revisions.BumpGlobalTx(tx)
		return err
	}); result.IsError() {
		rollback()
		return result
	}

	summary, err := s.roleAccessSummary(role, globalVersion+1)
	if err != nil {
		return s.internalFailure("build updated role response", err)
	}
	return response.ResultSuccess("Permissions assigned successfully", summary)
}

func (s *Service) persistMutation(entry *model.AuthorizationAuditLog, bump func(*gorm.DB) error) response.Result {
	tx := s.db.Begin()
	if tx.Error != nil {
		return s.dbFailure("begin authorization transaction", tx.Error)
	}
	if err := s.auditRepo.AppendTx(tx, entry); err != nil {
		_ = tx.Rollback().Error
		return s.dbFailure("append authorization audit", err)
	}
	if err := bump(tx); err != nil {
		_ = tx.Rollback().Error
		return s.dbFailure("bump authorization revision", err)
	}
	if err := tx.Commit().Error; err != nil {
		_ = tx.Rollback().Error
		return s.dbFailure("commit authorization transaction", err)
	}
	return response.ResultSuccess("Authorization mutation persisted", nil)
}

func (s *Service) roleManagerCount(ctx context.Context) (int, error) {
	roles, err := s.roleRepo.FindAll(ctx, nil, nil)
	if err != nil {
		return 0, err
	}
	managers := make(map[string]struct{})
	for _, role := range roles {
		permissions, err := s.policyManager.PermissionNamesForRole(role.ID.String(), authorization.OrgPlatform)
		if err != nil {
			return 0, err
		}
		if !containsString(permissions, "role:manage") {
			continue
		}
		users, err := s.policyManager.UsersForRole(role.ID.String(), authorization.OrgPlatform)
		if err != nil {
			return 0, err
		}
		for _, userID := range users {
			managers[userID] = struct{}{}
		}
	}
	return len(managers), nil
}

func (s *Service) subjectCanManageRoles(ctx context.Context, subject string) (bool, error) {
	if subject == authorization.SubjectAnonymous {
		return false, nil
	}
	err := s.authorizer.Enforce(ctx, authorization.Request{
		Subject: subject,
		Org:     authorization.OrgPlatform,
		Action:  authorization.ActionManage,
		Object:  authorization.ObjectRole,
		Context: authorization.CtxAny,
	})
	if errors.Is(err, authorization.ErrForbidden) {
		return false, nil
	}
	return err == nil, err
}

func newAuthorizationAudit(
	viewer authorization.ViewerContext,
	action string,
	targetType string,
	targetID uuid.UUID,
	targetName string,
	before common.JSONMap,
	after common.JSONMap,
) *model.AuthorizationAuditLog {
	entry := &model.AuthorizationAuditLog{
		ID: uuid.New(), Action: action, TargetType: targetType, TargetID: targetID,
		TargetNameSnapshot: targetName, Before: before, After: after, CreatedAt: time.Now().UTC(),
	}
	if viewer.User != nil {
		actorID := viewer.User.ID
		entry.ActorUserID = &actorID
		entry.ActorNameSnapshot = viewer.User.Name
		entry.ActorEmailSnapshot = viewer.User.Email
	}
	return entry
}

func authorizationConflict(code string, message string) response.Result {
	return response.ResultConflict(code, message)
}

func (s *Service) dbFailure(operation string, err error) response.Result {
	s.logMutationError(operation, err)
	return response.ResultErrDb(err)
}

func (s *Service) internalFailure(operation string, err error) response.Result {
	s.logMutationError(operation, err)
	return response.ResultErrInternal(err)
}

func (s *Service) logMutationError(operation string, err error) {
	if s.logger != nil {
		s.logger.Error("Authorization mutation failed", "operation", operation, "error", err)
	}
}

func uniqueUUIDs(values []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(values))
	result := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].String() < result[j].String() })
	return result
}

func stringifyUUIDs(values []uuid.UUID) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value.String())
	}
	return result
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func containsString(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func (s *Service) CreateRole(ctx context.Context, name string, description *string) response.Result {
	return s.locker.WithLock(ctx, func() response.Result {
		name, description, result := normalizeRoleMetadata(name, description)
		if result.IsError() {
			return result
		}
		viewer := authorization.ViewerFromContext(ctx)
		role := &model.Role{
			SqlModel:    common.SqlModel{ID: uuid.New()},
			Name:        name,
			Description: description,
		}
		entry := newAuthorizationAudit(viewer, "role.created", "role", role.ID, role.Name,
			common.JSONMap{}, roleMetadataSnapshot(role.Name, role.Description))
		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.roleRepo.CreateWithTransaction(tx, role); err != nil {
				return err
			}
			return s.auditRepo.AppendTx(tx, entry)
		})
		if err != nil {
			return s.dbFailure("create authorization role", err)
		}
		return response.ResultSuccess("Role created successfully", role)
	})
}

func (s *Service) UpdateRole(
	ctx context.Context,
	roleID uuid.UUID,
	name string,
	description *string,
	expectedVersion string,
) response.Result {
	return s.locker.WithLock(ctx, func() response.Result {
		name, description, result := normalizeRoleMetadata(name, description)
		if result.IsError() {
			return result
		}
		role, err := s.roleRepo.FindOne(ctx, []any{clause.Eq{Column: "id", Value: roleID}}, nil)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Role")
		}
		if err != nil {
			return s.dbFailure("find role for update", err)
		}
		globalVersion, err := s.revisions.CurrentGlobal(ctx)
		if err != nil {
			return s.dbFailure("read authorization revision", err)
		}
		if expectedVersion != "" && expectedVersion != fmt.Sprintf("g%d", globalVersion) {
			return authorizationConflict(codeAuthorizationStateChanged, "Authorization state changed; refresh and try again")
		}

		viewer := authorization.ViewerFromContext(ctx)
		before := roleMetadataSnapshot(role.Name, role.Description)
		after := roleMetadataSnapshot(name, description)
		entry := newAuthorizationAudit(viewer, "role.updated", "role", role.ID, name, before, after)
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.roleRepo.UpdateWithTransaction(tx,
				[]any{clause.Eq{Column: "id", Value: role.ID}},
				map[string]any{"name": name, "description": description},
			); err != nil {
				return err
			}
			if err := s.auditRepo.AppendTx(tx, entry); err != nil {
				return err
			}
			_, err := s.revisions.BumpGlobalTx(tx)
			return err
		})
		if err != nil {
			return s.dbFailure("update authorization role", err)
		}
		role.Name = name
		role.Description = description
		summary, err := s.roleAccessSummary(role, globalVersion+1)
		if err != nil {
			return s.internalFailure("build updated role response", err)
		}
		return response.ResultSuccess("Role updated successfully", summary)
	})
}

func (s *Service) DeleteRole(ctx context.Context, roleID uuid.UUID, expectedVersion string) response.Result {
	return s.locker.WithLock(ctx, func() response.Result {
		role, err := s.roleRepo.FindOne(ctx, []any{clause.Eq{Column: "id", Value: roleID}}, nil)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Role")
		}
		if err != nil {
			return s.dbFailure("find role for deletion", err)
		}
		globalVersion, err := s.revisions.CurrentGlobal(ctx)
		if err != nil {
			return s.dbFailure("read authorization revision", err)
		}
		if expectedVersion != "" && expectedVersion != fmt.Sprintf("g%d", globalVersion) {
			return authorizationConflict(codeAuthorizationStateChanged, "Authorization state changed; refresh and try again")
		}
		users, err := s.policyManager.UsersForRole(role.ID.String(), authorization.OrgPlatform)
		if err != nil {
			return s.internalFailure("read role assignments", err)
		}
		if len(users) > 0 {
			return authorizationConflict("ROLE_IN_USE", fmt.Sprintf("role is assigned to %d user(s)", len(users)))
		}
		permissions, err := s.policyManager.PermissionNamesForRole(role.ID.String(), authorization.OrgPlatform)
		if err != nil {
			return s.internalFailure("read role permissions", err)
		}
		if err := s.policyManager.RemoveRole(role.ID.String(), authorization.OrgPlatform); err != nil {
			return s.internalFailure("remove role policy", err)
		}
		rollbackPolicy := func() {
			if err := s.policyManager.ReplacePermissionsForRole(
				role.ID.String(), permissions, authorization.OrgPlatform,
			); err != nil {
				s.logMutationError("restore deleted role policy", err)
			}
		}

		viewer := authorization.ViewerFromContext(ctx)
		entry := newAuthorizationAudit(viewer, "role.deleted", "role", role.ID, role.Name,
			roleMetadataSnapshot(role.Name, role.Description), common.JSONMap{})
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.roleRepo.DeleteSoftWithTransaction(tx,
				[]any{clause.Eq{Column: "id", Value: role.ID}},
			); err != nil {
				return err
			}
			if err := s.auditRepo.AppendTx(tx, entry); err != nil {
				return err
			}
			_, err := s.revisions.BumpGlobalTx(tx)
			return err
		})
		if err != nil {
			rollbackPolicy()
			return s.dbFailure("delete authorization role", err)
		}
		return response.ResultSuccess("Role deleted successfully", nil)
	})
}

func normalizeRoleMetadata(name string, description *string) (string, *string, response.Result) {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 100 {
		return "", nil, response.ResultError("role name must contain between 2 and 100 characters")
	}
	if description != nil {
		trimmed := strings.TrimSpace(*description)
		if len(trimmed) > 1000 {
			return "", nil, response.ResultError("role description must not exceed 1000 characters")
		}
		description = &trimmed
	}
	return name, description, response.ResultSuccess("valid role metadata", nil)
}

func roleMetadataSnapshot(name string, description *string) common.JSONMap {
	var value any
	if description != nil {
		value = *description
	}
	return common.JSONMap{"name": name, "description": value}
}
