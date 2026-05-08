package readinghistoryroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ReadingHistoryRoute struct {
	*gin.Engine
	logger          *logger.Logger
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
	handler         *ReadingHistoryHandler
}

type ReadingHistoryRouteParams struct {
	fx.In

	Logger          *logger.Logger
	R               *gin.Engine
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
	Handler         *ReadingHistoryHandler
}

func NewReadingHistoryRoute(p ReadingHistoryRouteParams) *ReadingHistoryRoute {
	return &ReadingHistoryRoute{
		logger:          p.Logger,
		Engine:          p.R,
		authMiddleware:  p.AuthMiddleware,
		authzMiddleware: p.AuthzMiddleware,
		handler:         p.Handler,
	}
}

func (rhr *ReadingHistoryRoute) Setup() {
	rg := rhr.Group("/reading-histories", rhr.authMiddleware.RequireJwt)
	requireReadingHistoryCreate := authzmiddleware.Require(rhr.authzMiddleware, authorization.ActionCreate, authorization.ObjectReadingHistory)
	requireReadingHistoryRead := authzmiddleware.Require(rhr.authzMiddleware, authorization.ActionRead, authorization.ObjectReadingHistory, rhr.authzMiddleware.ReadingHistoryParam("id"))
	requireReadingHistoryUpdate := authzmiddleware.Require(rhr.authzMiddleware, authorization.ActionUpdate, authorization.ObjectReadingHistory, rhr.authzMiddleware.ReadingHistoryParam("id"))
	requireReadingHistoryDelete := authzmiddleware.Require(rhr.authzMiddleware, authorization.ActionDelete, authorization.ObjectReadingHistory, rhr.authzMiddleware.ReadingHistoryParam("id"))

	rg.GET("", rhr.handler.getReadingHistories)
	rg.GET("/:id", requireReadingHistoryRead, rhr.handler.getReadingHistory)
	rg.POST("", requireReadingHistoryCreate, rhr.handler.createReadingHistory)
	rg.PUT("/:id", requireReadingHistoryUpdate, rhr.handler.updateReadingHistory)
	rg.DELETE("/:id", requireReadingHistoryDelete, rhr.handler.deleteReadingHistory)
}
