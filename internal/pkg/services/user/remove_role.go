package userservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *UserService) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) response.Result {
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

	role, err := s.roleRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: roleID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Role")
		}
		s.logger.Error("Failed to find role", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.userRepo.RemoveRole(ctx, userID, role); err != nil {
		s.logger.Error("Failed to remove role from user", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Role removed successfully", nil)
}
