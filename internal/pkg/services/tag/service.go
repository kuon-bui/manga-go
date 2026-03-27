package tagservice

import (
	"manga-go/internal/pkg/logger"
	tagrepo "manga-go/internal/pkg/repo/tag"

	"go.uber.org/fx"
)

type TagService struct {
	logger  *logger.Logger
	tagRepo *tagrepo.TagRepo
}

type TagServiceParams struct {
	fx.In
	Logger  *logger.Logger
	TagRepo *tagrepo.TagRepo
}

func NewTagService(params TagServiceParams) *TagService {
	return &TagService{
		logger:  params.Logger,
		tagRepo: params.TagRepo,
	}
}
