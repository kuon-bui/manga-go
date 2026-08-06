package translationgrouproute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	slugmiddleware "manga-go/internal/app/middleware/slug"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type TranslationGroupRoute struct {
	*gin.Engine
	logger                  *logger.Logger
	authMiddleware          *authmiddleware.AuthMiddleware
	authzMiddleware         *authzmiddleware.AuthzMiddleware
	translationGroupHandler *TranslationGroupHandler
	slugMiddleware          *slugmiddleware.SlugMiddleware
}

type TranslationGroupRouteParams struct {
	fx.In

	R                       *gin.Engine
	Logger                  *logger.Logger
	AuthMiddleware          *authmiddleware.AuthMiddleware
	AuthzMiddleware         *authzmiddleware.AuthzMiddleware
	TranslationGroupHandler *TranslationGroupHandler
	SlugMiddleware          *slugmiddleware.SlugMiddleware
}

func NewTranslationGroupRoute(params TranslationGroupRouteParams) *TranslationGroupRoute {
	return &TranslationGroupRoute{
		logger:                  params.Logger,
		Engine:                  params.R,
		authMiddleware:          params.AuthMiddleware,
		authzMiddleware:         params.AuthzMiddleware,
		translationGroupHandler: params.TranslationGroupHandler,
		slugMiddleware:          params.SlugMiddleware,
	}
}

func (r *TranslationGroupRoute) Setup() {
	rg := r.Group("/translation-groups", r.authMiddleware.RequireJwt)
	requireGroupCreate := authzmiddleware.Require(r.authzMiddleware, authorization.ActionCreate, authorization.ObjectTranslationGroup)
	requireGroupUpdate := authzmiddleware.Require(r.authzMiddleware, authorization.ActionUpdate, authorization.ObjectTranslationGroup, r.authzMiddleware.TranslationGroup())
	requireGroupDelete := authzmiddleware.Require(r.authzMiddleware, authorization.ActionDelete, authorization.ObjectTranslationGroup, r.authzMiddleware.TranslationGroup())
	requireGroupManage := authzmiddleware.Require(r.authzMiddleware, authorization.ActionManage, authorization.ObjectTranslationGroup, r.authzMiddleware.TranslationGroup())

	rg.GET("", r.translationGroupHandler.getTranslationGroups)
	rg.POST("", requireGroupCreate, r.translationGroupHandler.createTranslationGroup)

	slugRg := r.Group("/translation-groups/:translationGroupSlug", r.authMiddleware.RequireJwt, r.slugMiddleware.ResolveTranslationGroupID)
	slugRg.GET("", r.translationGroupHandler.getTranslationGroup)
	slugRg.PUT("", requireGroupUpdate, r.translationGroupHandler.updateTranslationGroup)
	slugRg.DELETE("", requireGroupDelete, r.translationGroupHandler.deleteTranslationGroup)
	slugRg.PUT("/transfer-ownership", requireGroupManage, r.translationGroupHandler.transferOwnership)
	slugRg.GET("/members", r.translationGroupHandler.getMembers)
	slugRg.PUT("/logo", requireGroupUpdate, r.translationGroupHandler.updateLogo)
}
