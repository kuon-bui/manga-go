package roleservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *RoleService) GetRole(ctx context.Context, id uuid.UUID) response.Result {
	role, err := s.authAdmin.GetRole(ctx, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Role")
		}
		s.logger.Error("Failed to find role", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Role retrieved successfully", role)
}
