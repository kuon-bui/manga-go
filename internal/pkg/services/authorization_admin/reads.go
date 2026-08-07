package authorizationadmin

import (
	"context"
	"fmt"
	"sort"

	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/model"
	userrequest "manga-go/internal/pkg/request/user"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type ListUsersInput struct {
	Page   int
	Limit  int
	Search string
	RoleID *uuid.UUID
}

func (s *Service) ListUsers(ctx context.Context, input ListUsersInput) (*PagedUsers, error) {
	roleUserIDs, err := s.roleFilterUserIDs(input.RoleID)
	if err != nil {
		return nil, err
	}
	request := &userrequest.ListAuthorizationUsersRequest{
		Search: input.Search,
		RoleID: input.RoleID,
	}
	request.Page = input.Page
	request.Limit = input.Limit
	users, total, err := s.userRepo.ListAuthorizationUsers(ctx, request, roleUserIDs)
	if err != nil {
		return nil, err
	}

	userIDs := make([]uuid.UUID, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	globalVersion, userVersions, err := s.revisions.CurrentMany(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	roleModels, err := s.rolesForUsers(ctx, users)
	if err != nil {
		return nil, err
	}

	data := make([]AdminUserSummary, 0, len(users))
	for _, user := range users {
		roles, err := s.roleSummariesForUser(user.ID, roleModels)
		if err != nil {
			return nil, err
		}
		data = append(data, AdminUserSummary{
			ID:                   user.ID,
			Name:                 user.Name,
			Email:                user.Email,
			Roles:                roles,
			AuthorizationVersion: fmt.Sprintf("g%d:u%d", globalVersion, userVersions[user.ID]),
		})
	}
	return &PagedUsers{
		Data:  data,
		Total: total,
		Page:  request.Page,
		Limit: request.Limit,
	}, nil
}

func (s *Service) GetUser(ctx context.Context, userID uuid.UUID) (*AdminUserSummary, error) {
	user, err := s.userRepo.FindOne(ctx, []any{clause.Eq{Column: "id", Value: userID}}, nil)
	if err != nil {
		return nil, err
	}
	globalVersion, userVersion, err := s.revisions.Current(ctx, userID)
	if err != nil {
		return nil, err
	}
	roleModels, err := s.rolesForUsers(ctx, []*model.User{user})
	if err != nil {
		return nil, err
	}
	roles, err := s.roleSummariesForUser(user.ID, roleModels)
	if err != nil {
		return nil, err
	}
	return &AdminUserSummary{
		ID:                   user.ID,
		Name:                 user.Name,
		Email:                user.Email,
		Roles:                roles,
		AuthorizationVersion: fmt.Sprintf("g%d:u%d", globalVersion, userVersion),
	}, nil
}

func (s *Service) ListRoles(ctx context.Context) ([]RoleAccessSummary, error) {
	roles, err := s.roleRepo.FindAll(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	globalVersion, err := s.revisions.CurrentGlobal(ctx)
	if err != nil {
		return nil, err
	}
	summaries := make([]RoleAccessSummary, 0, len(roles))
	for _, role := range roles {
		summary, err := s.roleAccessSummary(role, globalVersion)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, *summary)
	}
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Name < summaries[j].Name
	})
	return summaries, nil
}

func (s *Service) GetRole(ctx context.Context, roleID uuid.UUID) (*RoleAccessSummary, error) {
	role, err := s.roleRepo.FindOne(ctx, []any{clause.Eq{Column: "id", Value: roleID}}, nil)
	if err != nil {
		return nil, err
	}
	globalVersion, err := s.revisions.CurrentGlobal(ctx)
	if err != nil {
		return nil, err
	}
	return s.roleAccessSummary(role, globalVersion)
}

func (s *Service) roleAccessSummary(role *model.Role, globalVersion uint64) (*RoleAccessSummary, error) {
	permissions, err := s.policyManager.PermissionNamesForRole(role.ID.String(), authorization.OrgPlatform)
	if err != nil {
		return nil, err
	}
	users, err := s.policyManager.UsersForRole(role.ID.String(), authorization.OrgPlatform)
	if err != nil {
		return nil, err
	}
	return &RoleAccessSummary{
		ID:                   role.ID,
		Name:                 role.Name,
		Description:          role.Description,
		Permissions:          permissions,
		AssignedUserCount:    len(users),
		AuthorizationVersion: fmt.Sprintf("g%d", globalVersion),
	}, nil
}

func (s *Service) roleFilterUserIDs(roleID *uuid.UUID) ([]uuid.UUID, error) {
	if roleID == nil {
		return nil, nil
	}
	rawIDs, err := s.policyManager.UsersForRole(roleID.String(), authorization.OrgPlatform)
	if err != nil {
		return nil, err
	}
	return parsePolicyUUIDs(rawIDs)
}

func (s *Service) rolesForUsers(ctx context.Context, users []*model.User) (map[uuid.UUID]*model.Role, error) {
	unique := make(map[uuid.UUID]struct{})
	for _, user := range users {
		rawIDs, err := s.policyManager.RolesForUser(user.ID.String(), authorization.OrgPlatform)
		if err != nil {
			return nil, err
		}
		ids, err := parsePolicyUUIDs(rawIDs)
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			unique[id] = struct{}{}
		}
	}
	if len(unique) == 0 {
		return map[uuid.UUID]*model.Role{}, nil
	}
	ids := make([]uuid.UUID, 0, len(unique))
	for id := range unique {
		ids = append(ids, id)
	}
	var roles []*model.Role
	if err := s.roleRepo.DB.WithContext(ctx).Where("id IN ?", ids).Find(&roles).Error; err != nil {
		return nil, err
	}
	byID := make(map[uuid.UUID]*model.Role, len(roles))
	for _, role := range roles {
		byID[role.ID] = role
	}
	return byID, nil
}

func (s *Service) roleSummariesForUser(userID uuid.UUID, models map[uuid.UUID]*model.Role) ([]RoleSummary, error) {
	rawIDs, err := s.policyManager.RolesForUser(userID.String(), authorization.OrgPlatform)
	if err != nil {
		return nil, err
	}
	ids, err := parsePolicyUUIDs(rawIDs)
	if err != nil {
		return nil, err
	}
	roles := make([]RoleSummary, 0, len(ids))
	for _, id := range ids {
		role, ok := models[id]
		if !ok {
			continue
		}
		roles = append(roles, RoleSummary{ID: role.ID, Name: role.Name, Description: role.Description})
	}
	sort.Slice(roles, func(i, j int) bool { return roles[i].Name < roles[j].Name })
	return roles, nil
}

func parsePolicyUUIDs(rawIDs []string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(rawIDs))
	for _, rawID := range rawIDs {
		id, err := uuid.Parse(rawID)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID %q in authorization policy: %w", rawID, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
