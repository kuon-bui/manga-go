package translationgrouproute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	casbinmiddleware "manga-go/internal/app/middleware/casbin"
	slugmiddleware "manga-go/internal/app/middleware/slug"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type TranslationGroupRoute struct {
	logger                  *logger.Logger
	r                       *gin.Engine
	authMiddleware          *authmiddleware.AuthMiddleware
	casbinMiddleware        *casbinmiddleware.CasbinMiddleware
	slugMiddleware          *slugmiddleware.SlugMiddleware
	translationGroupHandler *TranslationGroupHandler
}

type TranslationGroupRouteParams struct {
	fx.In

	R                       *gin.Engine
	Logger                  *logger.Logger
	AuthMiddleware          *authmiddleware.AuthMiddleware
	CasbinMiddleware        *casbinmiddleware.CasbinMiddleware
	SlugMiddleware          *slugmiddleware.SlugMiddleware
	TranslationGroupHandler *TranslationGroupHandler
}

func NewTranslationGroupRoute(params TranslationGroupRouteParams) *TranslationGroupRoute {
	return &TranslationGroupRoute{
		logger:                  params.Logger,
		r:                       params.R,
		authMiddleware:          params.AuthMiddleware,
		casbinMiddleware:        params.CasbinMiddleware,
		slugMiddleware:          params.SlugMiddleware,
		translationGroupHandler: params.TranslationGroupHandler,
	}
}

func (r *TranslationGroupRoute) Setup() {
	rg := r.r.Group("/translation-groups")

	// Public read endpoints
	rg.GET("/", r.translationGroupHandler.getTranslationGroups)
	rg.GET("/:slug", r.translationGroupHandler.getTranslationGroup)

	// Authenticated endpoints
	auth := rg.Group("", r.authMiddleware.RequireJwt)

	// Join a translation group (non-group users)
	auth.POST("/join", r.translationGroupHandler.joinTranslationGroup)

	// Create a new translation group
	auth.POST("/", r.translationGroupHandler.createTranslationGroup)

	// Group owner management – resolves group ID into context for Casbin checks
	ownerRoutes := auth.Group("/:slug", r.slugMiddleware.ResolveTranslationGroupID)
	ownerRoutes.PUT("/", r.translationGroupHandler.updateTranslationGroup)
	ownerRoutes.DELETE("/", r.translationGroupHandler.deleteTranslationGroup)
	ownerRoutes.PUT("/transfer-ownership", r.translationGroupHandler.transferOwnership)
	ownerRoutes.POST("/kick", r.translationGroupHandler.kickMember)
	ownerRoutes.POST("/grant-permission", r.translationGroupHandler.grantPermission)
}
