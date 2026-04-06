package tagservice

import (
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/redis"
	tagrepo "manga-go/internal/pkg/repo/tag"

	"go.uber.org/fx"
)

type TagService struct {
	logger  *logger.Logger
	tagRepo *tagrepo.TagRepo
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
