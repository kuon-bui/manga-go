package permissionservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *PermissionService) DeletePermission(ctx context.Context, id uuid.UUID) response.Result {
	permission, err := s.permissionRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Permission")
		}
		s.logger.Error("Failed to find permission for deletion", "error", err)
		return response.ResultErrDb(err)
	}

	if err := s.permissionRepo.DeleteSoft(ctx, []any{
		clause.Eq{Column: "id", Value: id},
	}); err != nil {
		s.logger.Error("Failed to delete permission", "error", err)
		return response.ResultErrDb(err)
	}

	if s.policyManager != nil {
		if err := s.policyManager.RemovePermission(permission.ID.String(), authorization.OrgPlatform); err != nil {
			s.logger.Error("Failed to remove authorization policy", "error", err)
			return response.ResultErrInternal(err)
		}
	}

	return response.ResultSuccess("Permission deleted successfully", nil)
}
