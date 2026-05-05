package roleroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RoleRoute struct {
	*gin.Engine
	logger         *logger.Logger
	authMiddleware *authmiddleware.AuthMiddleware
	roleHandler    *RoleHandler
}

type RoleRouteParams struct {
	fx.In

	R              *gin.Engine
	Logger         *logger.Logger
	RoleHandler    *RoleHandler
	AuthMiddleware *authmiddleware.AuthMiddleware
}

func NewRoleRoute(params RoleRouteParams) *RoleRoute {
	return &RoleRoute{
		logger:         params.Logger,
		Engine:         params.R,
		authMiddleware: params.AuthMiddleware,
		roleHandler:    params.RoleHandler,
	}
}

func (rr *RoleRoute) Setup() {
	rg := rr.Group("/roles", rr.authMiddleware.RequireJwt)

	rg.GET("", rr.roleHandler.getRoles)
	rg.GET("/all", rr.roleHandler.getAllRoles)
	rg.GET("/:id", rr.roleHandler.getRole)
	rg.POST("", rr.roleHandler.createRole)
	rg.PUT("/:id", rr.roleHandler.updateRole)
	rg.DELETE("/:id", rr.roleHandler.deleteRole)
	rg.POST("/:id/permissions", rr.roleHandler.assignRolePermission)
	rg.DELETE("/:id/permissions/:permissionId", rr.roleHandler.removeRolePermission)
}
