package genreroute

import (
	genreservice "manga-go/internal/pkg/services/genre"

	"go.uber.org/fx"
)

type GenreHandler struct {
	genreService *genreservice.GenreService
}

type GenreHandlerParams struct {
	fx.In

	GenreService *genreservice.GenreService
}

func NewGenreHandler(p GenreHandlerParams) *GenreHandler {
	return &GenreHandler{
		genreService: p.GenreService,
	}
}
