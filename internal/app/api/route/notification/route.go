package notificationroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type NotificationRoute struct {
	*gin.Engine
	logger         *logger.Logger
	authMiddleware *authmiddleware.AuthMiddleware
	handler        *NotificationHandler
}

type NotificationRouteParams struct {
	fx.In
	Logger         *logger.Logger
	R              *gin.Engine
	AuthMiddleware *authmiddleware.AuthMiddleware
	Handler        *NotificationHandler
}

func NewNotificationRoute(p NotificationRouteParams) *NotificationRoute {
	return &NotificationRoute{
		logger:         p.Logger,
		Engine:         p.R,
		authMiddleware: p.AuthMiddleware,
		handler:        p.Handler,
	}
}

func (r *NotificationRoute) Setup() {
	rg := r.Group("/notifications", r.authMiddleware.RequireJwt)
	rg.GET("", r.handler.getNotifications)
	rg.GET("/stream", r.handler.streamNotifications)
	rg.PATCH("/:id/seen", r.handler.markNotificationSeen)
	rg.PATCH("/:id/read", r.handler.markNotificationRead)
	rg.PATCH("/read-all", r.handler.markAllNotificationsRead)
}
