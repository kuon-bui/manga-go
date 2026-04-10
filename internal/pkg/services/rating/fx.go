package ratingservice

import "go.uber.org/fx"

var Module = fx.Module(
	"rating-service",
	fx.Provide(NewRatingService),
)
