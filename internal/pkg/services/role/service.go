package roleservice

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	rolerepo "manga-go/internal/pkg/repo/role"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

// RoleRepository defines the data access interface for Role.
type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Role, error)
	Update(ctx context.Context, conditions []any, data map[string]any) error
	DeleteSoft(ctx context.Context, conditions []any) error
	FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Role, int64, error)
	FindAll(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*model.Role, error)
	AssignPermissions(ctx context.Context, roleID uuid.UUID, perms []*model.Permission) error
	RemovePermission(ctx context.Context, roleID uuid.UUID, perm *model.Permission) error
}

// PermissionRepository defines the data access interface for Permission (used by RoleService).
type PermissionRepository interface {
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Permission, error)
}

type RoleService struct {
	logger         *logger.Logger
	roleRepo       RoleRepository
	permissionRepo PermissionRepository
}

type RoleServiceParams struct {
	fx.In
	Logger         *logger.Logger
	RoleRepo       *rolerepo.RoleRepo
	PermissionRepo *permissionrepo.PermissionRepo
}

func NewRoleService(params RoleServiceParams) *RoleService {
	return &RoleService{
		logger:         params.Logger,
		roleRepo:       params.RoleRepo,
		permissionRepo: params.PermissionRepo,
	}
}

// NewRoleServiceWithRepos creates a RoleService with explicit repository interfaces,
// useful for unit testing.
func NewRoleServiceWithRepos(l *logger.Logger, roleRepo RoleRepository, permRepo PermissionRepository) *RoleService {
	return &RoleService{
		logger:         l,
		roleRepo:       roleRepo,
		permissionRepo: permRepo,
	}
}
