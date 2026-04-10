package ratingroute

import (
	ratingservice "manga-go/internal/pkg/services/rating"

	"go.uber.org/fx"
)

type RatingHandler struct {
	ratingService *ratingservice.RatingService
}

type RatingHandlerParams struct {
	fx.In

	RatingService *ratingservice.RatingService
}

func NewRatingHandler(p RatingHandlerParams) *RatingHandler {
	return &RatingHandler{
		ratingService: p.RatingService,
	}
}
