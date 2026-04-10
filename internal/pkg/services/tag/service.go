package tagservice

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/redis"
	tagrepo "manga-go/internal/pkg/repo/tag"

	"go.uber.org/fx"
)

// TagRepository defines the data access interface for Tag.
type TagRepository interface {
	Create(ctx context.Context, tag *model.Tag) error
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Tag, error)
	DeleteSoft(ctx context.Context, conditions []any) error
	FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Tag, int64, error)
	FindAll(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) ([]*model.Tag, error)
}

type TagService struct {
	logger  *logger.Logger
	tagRepo TagRepository
	rds     *redis.Redis
}

type TagServiceParams struct {
	fx.In
	Logger  *logger.Logger
	TagRepo *tagrepo.TagRepo
	Redis   *redis.Redis
}

func NewTagService(params TagServiceParams) *TagService {
	return &TagService{
		logger:  params.Logger,
		tagRepo: params.TagRepo,
		rds:     params.Redis,
	}
}

// NewTagServiceWithRepo creates a TagService with an explicit repository,
// useful for unit testing.
func NewTagServiceWithRepo(l *logger.Logger, repo TagRepository) *TagService {
	return &TagService{
		logger:  l,
		tagRepo: repo,
	}
}
