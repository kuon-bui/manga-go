package ratingroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	slugmiddleware "manga-go/internal/app/middleware/slug"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RatingRoute struct {
	*gin.Engine
	logger          *logger.Logger
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
	ratingHandler   *RatingHandler
	slugMiddleware  *slugmiddleware.SlugMiddleware
}

type RatingRouteParams struct {
	fx.In

	R               *gin.Engine
	Logger          *logger.Logger
	RatingHandler   *RatingHandler
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
	SlugMiddleware  *slugmiddleware.SlugMiddleware
}

func NewRatingRoute(params RatingRouteParams) *RatingRoute {
	return &RatingRoute{
		logger:          params.Logger,
		Engine:          params.R,
		authMiddleware:  params.AuthMiddleware,
		authzMiddleware: params.AuthzMiddleware,
		ratingHandler:   params.RatingHandler,
		slugMiddleware:  params.SlugMiddleware,
	}
}

func (rr *RatingRoute) Setup() {
	rg := rr.Group("/ratings/comics/:comicSlug", rr.authMiddleware.RequireJwt, rr.slugMiddleware.ResolveComicID)
	requireRatingCreate := authzmiddleware.Require(rr.authzMiddleware, authorization.ActionCreate, authorization.ObjectRating)
	requireRatingUpdate := authzmiddleware.Require(rr.authzMiddleware, authorization.ActionUpdate, authorization.ObjectRating, rr.authzMiddleware.RatingParam("id"))
	requireRatingDelete := authzmiddleware.Require(rr.authzMiddleware, authorization.ActionDelete, authorization.ObjectRating, rr.authzMiddleware.RatingParam("id"))

	rg.GET("", rr.ratingHandler.getRatings)
	rg.GET("average", rr.ratingHandler.getAverageRating)
	rg.POST("", requireRatingCreate, rr.ratingHandler.createRating)
	rg.PUT("/:id", requireRatingUpdate, rr.ratingHandler.updateRating)
	rg.DELETE("/:id", requireRatingDelete, rr.ratingHandler.deleteRating)
}
