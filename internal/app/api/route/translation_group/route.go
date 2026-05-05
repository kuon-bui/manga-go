package translationgrouproute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	slugmiddleware "manga-go/internal/app/middleware/slug"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type TranslationGroupRoute struct {
	*gin.Engine
	logger                  *logger.Logger
	authMiddleware          *authmiddleware.AuthMiddleware
	translationGroupHandler *TranslationGroupHandler
	slugMiddleware          *slugmiddleware.SlugMiddleware
}

type TranslationGroupRouteParams struct {
	fx.In

	R                       *gin.Engine
	Logger                  *logger.Logger
	AuthMiddleware          *authmiddleware.AuthMiddleware
	TranslationGroupHandler *TranslationGroupHandler
	SlugMiddleware          *slugmiddleware.SlugMiddleware
}

func NewTranslationGroupRoute(params TranslationGroupRouteParams) *TranslationGroupRoute {
	return &TranslationGroupRoute{
		logger:                  params.Logger,
		Engine:                  params.R,
		authMiddleware:          params.AuthMiddleware,
		translationGroupHandler: params.TranslationGroupHandler,
		slugMiddleware:          params.SlugMiddleware,
	}
}

func (r *TranslationGroupRoute) Setup() {
	rg := r.Group("/translation-groups", r.authMiddleware.RequireJwt)

	rg.GET("", r.translationGroupHandler.getTranslationGroups)
	rg.POST("", r.translationGroupHandler.createTranslationGroup)

	slugRg := r.Group("/translation-groups/:translationGroupSlug", r.authMiddleware.RequireJwt, r.slugMiddleware.ResolveTranslationGroupID)
	slugRg.GET("", r.translationGroupHandler.getTranslationGroup)
	slugRg.PUT("", r.translationGroupHandler.updateTranslationGroup)
	slugRg.DELETE("", r.translationGroupHandler.deleteTranslationGroup)
	slugRg.PUT("/transfer-ownership", r.translationGroupHandler.transferOwnership)
	slugRg.GET("/members", r.translationGroupHandler.getMembers)
	slugRg.PUT("/logo", r.translationGroupHandler.updateLogo)
}
