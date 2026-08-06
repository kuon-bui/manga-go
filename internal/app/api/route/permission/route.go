package permissionroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type PermissionRoute struct {
	*gin.Engine
	logger            *logger.Logger
	authMiddleware    *authmiddleware.AuthMiddleware
	authzMiddleware   *authzmiddleware.AuthzMiddleware
	permissionHandler *PermissionHandler
}

type PermissionRouteParams struct {
	fx.In

	R                 *gin.Engine
	Logger            *logger.Logger
	PermissionHandler *PermissionHandler
	AuthMiddleware    *authmiddleware.AuthMiddleware
	AuthzMiddleware   *authzmiddleware.AuthzMiddleware
}

func NewPermissionRoute(params PermissionRouteParams) *PermissionRoute {
	return &PermissionRoute{
		logger:            params.Logger,
		Engine:            params.R,
		authMiddleware:    params.AuthMiddleware,
		authzMiddleware:   params.AuthzMiddleware,
		permissionHandler: params.PermissionHandler,
	}
}

func (pr *PermissionRoute) Setup() {
	rg := pr.Group("/permissions", pr.authMiddleware.RequireJwt)
	requirePermissionManage := authzmiddleware.Require(pr.authzMiddleware, authorization.ActionManage, authorization.ObjectPermission)

	rg.GET("", requirePermissionManage, pr.permissionHandler.getPermissions)
	rg.GET("/all", requirePermissionManage, pr.permissionHandler.getAllPermissions)
	rg.POST("", requirePermissionManage, pr.permissionHandler.createPermission)
	rg.PUT("/:id", requirePermissionManage, pr.permissionHandler.updatePermission)
	rg.DELETE("/:id", requirePermissionManage, pr.permissionHandler.deletePermission)
}
