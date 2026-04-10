package ratingroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	slugmiddleware "manga-go/internal/app/middleware/slug"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RatingRoute struct {
	logger         *logger.Logger
	r              *gin.Engine
	authMiddleware *authmiddleware.AuthMiddleware
	ratingHandler  *RatingHandler
	slugMiddleware *slugmiddleware.SlugMiddleware
}

type RatingRouteParams struct {
	fx.In

	R              *gin.Engine
	Logger         *logger.Logger
	RatingHandler  *RatingHandler
	AuthMiddleware *authmiddleware.AuthMiddleware
	SlugMiddleware *slugmiddleware.SlugMiddleware
}

func NewRatingRoute(params RatingRouteParams) *RatingRoute {
	return &RatingRoute{
		logger:         params.Logger,
		r:              params.R,
		authMiddleware: params.AuthMiddleware,
		ratingHandler:  params.RatingHandler,
		slugMiddleware: params.SlugMiddleware,
	}
}

func (rr *RatingRoute) Setup() {
	rg := rr.r.Group("/ratings/comics/:comicSlug", rr.authMiddleware.RequireJwt, rr.slugMiddleware.ResolveComicID)

	rg.GET("", rr.ratingHandler.getRatings)
	rg.GET("average", rr.ratingHandler.getAverageRating)
	rg.POST("", rr.ratingHandler.createRating)
	rg.PUT("/:id", rr.ratingHandler.updateRating)
	rg.DELETE("/:id", rr.ratingHandler.deleteRating)
}
