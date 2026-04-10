package ratingrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"rating-repo",
	fx.Provide(NewRatingRepo),
)
