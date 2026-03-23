package genreroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"genre-route",
	common.ProvideAsRoute(NewGenreRoute),
	fx.Provide(NewGenreHandler),
)
