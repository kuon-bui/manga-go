package userservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

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

	roles, err := s.userRepo.GetRoles(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user roles", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("User roles retrieved successfully", roles)
}
