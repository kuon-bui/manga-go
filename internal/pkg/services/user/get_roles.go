package userservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *UserService) GetUserRoles(ctx context.Context, userID uuid.UUID) response.Result {
	_, err := s.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: userID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("User")
		}
		s.logger.Error("Failed to find user", "error", err)
		return response.ResultErrDb(err)
	}

	// Who holds which role is recorded in the policy engine; the roles table
	// only carries the metadata to display alongside it.
	roleIDs, err := s.policyManager.RolesForUser(userID.String(), authorization.OrgPlatform)
	if err != nil {
		s.logger.Error("Failed to read user roles", "error", err)
		return response.ResultErrInternal(err)
	}

	roles := make([]*model.Role, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		id, err := uuid.Parse(roleID)
		if err != nil {
			// A non-uuid subject is a built-in role such as group_member, which
			// has no row to describe it.
			continue
		}

		role, err := s.roleRepo.FindOne(ctx, []any{
			clause.Eq{Column: "id", Value: id},
		}, nil)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			s.logger.Error("Failed to find role", "error", err)
			return response.ResultErrDb(err)
		}
		roles = append(roles, role)
	}

	return response.ResultSuccess("User roles retrieved successfully", roles)
}
