package permissionroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type PermissionRoute struct {
	*gin.Engine
	logger            *logger.Logger
	authMiddleware    *authmiddleware.AuthMiddleware
	permissionHandler *PermissionHandler
}

type PermissionRouteParams struct {
	fx.In

	R                 *gin.Engine
	Logger            *logger.Logger
	PermissionHandler *PermissionHandler
	AuthMiddleware    *authmiddleware.AuthMiddleware
}

func NewPermissionRoute(params PermissionRouteParams) *PermissionRoute {
	return &PermissionRoute{
		logger:            params.Logger,
		Engine:            params.R,
		authMiddleware:    params.AuthMiddleware,
		permissionHandler: params.PermissionHandler,
	}
}

func (pr *PermissionRoute) Setup() {
	rg := pr.Group("/permissions", pr.authMiddleware.RequireJwt)

	rg.GET("", pr.permissionHandler.getPermissions)
	rg.GET("/all", pr.permissionHandler.getAllPermissions)
	rg.POST("", pr.permissionHandler.createPermission)
	rg.PUT("/:id", pr.permissionHandler.updatePermission)
	rg.DELETE("/:id", pr.permissionHandler.deletePermission)
}
