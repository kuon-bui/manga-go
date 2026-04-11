package permissionservice

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	permissionrepo "manga-go/internal/pkg/repo/permission"

	"go.uber.org/fx"
)

// PermissionRepository defines the data access interface for Permission.
type PermissionRepository interface {
	Create(ctx context.Context, permission *model.Permission) error
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Permission, error)
	Update(ctx context.Context, conditions []any, data map[string]any) error
	DeleteSoft(ctx context.Context, conditions []any) error
	FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Permission, int64, error)
	FindAll(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*model.Permission, error)
}

// Ensure *permissionrepo.PermissionRepo satisfies PermissionRepository.
var _ PermissionRepository = (*permissionrepo.PermissionRepo)(nil)

type PermissionService struct {
	logger         *logger.Logger
	permissionRepo PermissionRepository
}

type PermissionServiceParams struct {
	fx.In
	Logger         *logger.Logger
	PermissionRepo *permissionrepo.PermissionRepo
}

func NewPermissionService(params PermissionServiceParams) *PermissionService {
	return &PermissionService{
		logger:         params.Logger,
		permissionRepo: params.PermissionRepo,
	}
}

// NewPermissionServiceWithRepo creates a PermissionService with an explicit repository,
// useful for unit testing.
func NewPermissionServiceWithRepo(l *logger.Logger, repo PermissionRepository) *PermissionService {
	return &PermissionService{
		logger:         l,
		permissionRepo: repo,
	}
}
