package genreservice

import (
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/redis"
	genrerepo "manga-go/internal/pkg/repo/genre"

	"go.uber.org/fx"
)

type GenreService struct {
	logger    *logger.Logger
	genreRepo *genrerepo.GenreRepo
	rds       *redis.Redis
}

type GenreServiceParams struct {
	fx.In
	Logger    *logger.Logger
	GenreRepo *genrerepo.GenreRepo
	Redis     *redis.Redis
}

func NewGenreService(params GenreServiceParams) *GenreService {
	return &GenreService{
		logger:    params.Logger,
		genreRepo: params.GenreRepo,
		rds:       params.Redis,
	}
}
