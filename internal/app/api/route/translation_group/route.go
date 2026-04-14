package translationgrouproute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type TranslationGroupRoute struct {
	logger                  *logger.Logger
	r                       *gin.Engine
	authMiddleware          *authmiddleware.AuthMiddleware
	translationGroupHandler *TranslationGroupHandler
}

type TranslationGroupRouteParams struct {
	fx.In

	R                       *gin.Engine
	Logger                  *logger.Logger
	AuthMiddleware          *authmiddleware.AuthMiddleware
	TranslationGroupHandler *TranslationGroupHandler
}

func NewTranslationGroupRoute(params TranslationGroupRouteParams) *TranslationGroupRoute {
	return &TranslationGroupRoute{
		logger:                  params.Logger,
		r:                       params.R,
		authMiddleware:          params.AuthMiddleware,
		translationGroupHandler: params.TranslationGroupHandler,
	}
}

func (r *TranslationGroupRoute) Setup() {
	rg := r.r.Group("/translation-groups", r.authMiddleware.RequireJwt)

	rg.GET("", r.translationGroupHandler.getTranslationGroups)
	rg.GET("/:slug", r.translationGroupHandler.getTranslationGroup)
	rg.POST("", r.translationGroupHandler.createTranslationGroup)
	rg.PUT("/:slug", r.translationGroupHandler.updateTranslationGroup)
	rg.DELETE("/:slug", r.translationGroupHandler.deleteTranslationGroup)
	rg.PUT("/:slug/transfer-ownership", r.translationGroupHandler.transferOwnership)
}
