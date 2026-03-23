package genreservice

import (
	"manga-go/internal/pkg/logger"
	genrerepo "manga-go/internal/pkg/repo/genre"

	"go.uber.org/fx"
)

type GenreService struct {
	logger    *logger.Logger
	genreRepo *genrerepo.GenreRepo
}

type GenreServiceParams struct {
	fx.In
	Logger    *logger.Logger
	GenreRepo *genrerepo.GenreRepo
}

func NewGenreService(params GenreServiceParams) *GenreService {
	return &GenreService{
		logger:    params.Logger,
		genreRepo: params.GenreRepo,
	}
}
