package genrerepo

import "go.uber.org/fx"

var Module = fx.Module(
	"genre-repo",
	fx.Provide(NewGenreRepo),
)
