package comicservice

import (
	"manga-go/internal/pkg/logger"
	comicrepo "manga-go/internal/pkg/repo/comic"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ComicService struct {
	logger    *logger.Logger
	comicRepo *comicrepo.ComicRepo
	db        *gorm.DB
}

type ComicServiceParams struct {
	fx.In
	Logger    *logger.Logger
	ComicRepo *comicrepo.ComicRepo
	DB        *gorm.DB
}

func NewComicService(params ComicServiceParams) *ComicService {
	return &ComicService{
		logger:    params.Logger,
		comicRepo: params.ComicRepo,
		db:        params.DB,
	}
}
