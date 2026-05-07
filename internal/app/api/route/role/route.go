package roleroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RoleRoute struct {
	*gin.Engine
	logger          *logger.Logger
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
	roleHandler     *RoleHandler
}

type RoleRouteParams struct {
	fx.In

	R               *gin.Engine
	Logger          *logger.Logger
	RoleHandler     *RoleHandler
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
}

func NewRoleRoute(params RoleRouteParams) *RoleRoute {
	return &RoleRoute{
		logger:          params.Logger,
		Engine:          params.R,
		authMiddleware:  params.AuthMiddleware,
		authzMiddleware: params.AuthzMiddleware,
		roleHandler:     params.RoleHandler,
	}
}

func (rr *RoleRoute) Setup() {
	rg := rr.Group("/roles", rr.authMiddleware.RequireJwt)
	requireRoleManage := authzmiddleware.Require(rr.authzMiddleware, authorization.ActionManage, authorization.ObjectRole)

	rg.GET("", requireRoleManage, rr.roleHandler.getRoles)
	rg.GET("/all", requireRoleManage, rr.roleHandler.getAllRoles)
	rg.GET("/:id", requireRoleManage, rr.roleHandler.getRole)
	rg.POST("", requireRoleManage, rr.roleHandler.createRole)
	rg.PUT("/:id", requireRoleManage, rr.roleHandler.updateRole)
	rg.DELETE("/:id", requireRoleManage, rr.roleHandler.deleteRole)
	rg.POST("/:id/permissions", requireRoleManage, rr.roleHandler.assignRolePermission)
	rg.DELETE("/:id/permissions/:permissionId", requireRoleManage, rr.roleHandler.removeRolePermission)
}
