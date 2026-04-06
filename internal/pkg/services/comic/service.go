package comicservice

import (
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/redis"
	comicrepo "manga-go/internal/pkg/repo/comic"
	genrerepo "manga-go/internal/pkg/repo/genre"
	tagrepo "manga-go/internal/pkg/repo/tag"

	"go.uber.org/fx"
)

type ComicService struct {
	logger    *logger.Logger
	comicRepo *comicrepo.ComicRepo
	genreRepo *genrerepo.GenreRepo
	tagRepo   *tagrepo.TagRepo
	rds       *redis.Redis
}

type ComicServiceParams struct {
	fx.In
	Logger    *logger.Logger
	ComicRepo *comicrepo.ComicRepo
	GenreRepo *genrerepo.GenreRepo
	TagRepo   *tagrepo.TagRepo
	Redis     *redis.Redis
}

func NewComicService(params ComicServiceParams) *ComicService {
	return &ComicService{
		logger:    params.Logger,
		comicRepo: params.ComicRepo,
		genreRepo: params.GenreRepo,
		tagRepo:   params.TagRepo,
		rds:       params.Redis,
	}
}
