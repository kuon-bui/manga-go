package translationgroupservice

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/redis"
	translationgrouprepo "manga-go/internal/pkg/repo/translation_group"
	userrepo "manga-go/internal/pkg/repo/user"

	"go.uber.org/fx"
)

// TranslationGroupRepository defines the data access interface for TranslationGroup.
type TranslationGroupRepository interface {
	Create(ctx context.Context, group *model.TranslationGroup) error
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.TranslationGroup, error)
	Update(ctx context.Context, conditions []any, data map[string]any) error
	DeleteSoft(ctx context.Context, conditions []any) error
	FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.TranslationGroup, int64, error)
}

// UserRepository defines the subset of UserRepo used by TranslationGroupService.
type UserRepository interface {
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.User, error)
	Update(ctx context.Context, conditions []any, data map[string]any) error
}

type TranslationGroupService struct {
	logger               *logger.Logger
	translationGroupRepo TranslationGroupRepository
	userRepo             UserRepository
	rds                  *redis.Redis
}

type TranslationGroupServiceParams struct {
	fx.In
	Logger               *logger.Logger
	TranslationGroupRepo *translationgrouprepo.TranslationGroupRepo
	UserRepo             *userrepo.UserRepository
	Redis                *redis.Redis
}

func NewTranslationGroupService(params TranslationGroupServiceParams) *TranslationGroupService {
	return &TranslationGroupService{
		logger:               params.Logger,
		translationGroupRepo: params.TranslationGroupRepo,
		userRepo:             params.UserRepo,
		rds:                  params.Redis,
	}
}

// NewTranslationGroupServiceWithRepos creates a TranslationGroupService with explicit repository interfaces,
// useful for unit testing.
func NewTranslationGroupServiceWithRepos(l *logger.Logger, tgRepo TranslationGroupRepository, userRepo UserRepository) *TranslationGroupService {
	return &TranslationGroupService{
		logger:               l,
		translationGroupRepo: tgRepo,
		userRepo:             userRepo,
	}
}
