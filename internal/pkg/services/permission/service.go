package permissionservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"go.uber.org/fx"
)

// PermissionService serves the permission catalog. Permissions are defined in
// code rather than stored as rows, so there is nothing here to create, update
// or delete — only the vocabulary an administrator can grant to a role.
type PermissionService struct {
	logger *logger.Logger
}

type PermissionServiceParams struct {
	fx.In
	Logger *logger.Logger
}

func NewPermissionService(params PermissionServiceParams) *PermissionService {
	return &PermissionService{logger: params.Logger}
}

func (s *PermissionService) ListPermissions(ctx context.Context) response.Result {
	return response.ResultSuccess("Permissions retrieved successfully", authorization.Catalog())
}
