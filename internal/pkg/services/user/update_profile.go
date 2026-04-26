package userservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	userrequest "manga-go/internal/pkg/request/user"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *userrequest.UpdateUserProfileRequest) response.Result {
	if req.Name == nil && req.Avatar == nil {
		return response.ResultError("At least one field is required")
	}

	updateData := map[string]any{}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return response.ResultError("Name cannot be empty")
		}
		updateData["name"] = name
	}

	if req.Avatar != nil {
		avatar := strings.TrimSpace(*req.Avatar)
		if avatar == "" {
			updateData["avatar"] = nil
		} else {
			updateData["avatar"] = avatar
		}
	}

	err := s.userRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: userID},
	}, updateData)
	if err != nil {
		s.logger.Error("Failed to update user profile", "error", err)
		return response.ResultErrDb(err)
	}

	updatedUser, err := s.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: userID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("User")
		}
		s.logger.Error("Failed to get updated user profile", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("User profile updated successfully", updatedUser)
}
