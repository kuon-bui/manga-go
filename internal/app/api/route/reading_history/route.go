package readinghistoryroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ReadingHistoryRoute struct {
	*gin.Engine
	logger         *logger.Logger
	authMiddleware *authmiddleware.AuthMiddleware
	handler        *ReadingHistoryHandler
}

type ReadingHistoryRouteParams struct {
	fx.In

	Logger         *logger.Logger
	R              *gin.Engine
	AuthMiddleware *authmiddleware.AuthMiddleware
	Handler        *ReadingHistoryHandler
}

func NewReadingHistoryRoute(p ReadingHistoryRouteParams) *ReadingHistoryRoute {
	return &ReadingHistoryRoute{
		logger:         p.Logger,
		Engine:         p.R,
		authMiddleware: p.AuthMiddleware,
		handler:        p.Handler,
	}
}

func (rhr *ReadingHistoryRoute) Setup() {
	rg := rhr.Group("/reading-histories", rhr.authMiddleware.RequireJwt)

	rg.GET("", rhr.handler.getReadingHistories)
	rg.GET("/:id", rhr.handler.getReadingHistory)
	rg.POST("", rhr.handler.createReadingHistory)
	rg.PUT("/:id", rhr.handler.updateReadingHistory)
	rg.DELETE("/:id", rhr.handler.deleteReadingHistory)
}
