package genreservice

import "go.uber.org/fx"

var Module = fx.Module(
	"genre-service",
	fx.Provide(NewGenreService),
)
