package ratingroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"rating-route",
	common.ProvideAsRoute(NewRatingRoute),
	fx.Provide(NewRatingHandler),
)
