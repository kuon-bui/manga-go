package userservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AssignRoles replaces the user's platform roles with exactly roleIDs. Each id
// is checked against the roles table so the policy engine only ever references
// roles that exist.
func (s *UserService) AssignRoles(
	ctx context.Context,
	userID uuid.UUID,
	roleIDs []uuid.UUID,
	expectedVersion ...string,
) response.Result {
	if s.authAdmin != nil && s.authAdmin.MutationReady() {
		return s.authAdmin.ReplaceUserRoles(ctx, userID, roleIDs, firstVersion(expectedVersion))
	}
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

	ids := make([]string, 0, len(roleIDs))
	for _, id := range roleIDs {
		role, err := s.roleRepo.FindOne(ctx, []any{
			clause.Eq{Column: "id", Value: id},
		}, nil)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.ResultNotFound("Role")
			}
			s.logger.Error("Failed to find role", "error", err)
			return response.ResultErrDb(err)
		}
		ids = append(ids, role.ID.String())
	}

	if err := s.policyManager.ReplaceRolesForUser(userID.String(), ids, authorization.OrgPlatform); err != nil {
		s.logger.Error("Failed to update authorization policy", "error", err)
		return response.ResultErrInternal(err)
	}

	return response.ResultSuccess("Roles assigned successfully", nil)
}

func firstVersion(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
