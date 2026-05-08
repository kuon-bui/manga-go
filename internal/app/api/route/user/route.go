package userroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type UserRoute struct {
	*gin.Engine
	logger          *logger.Logger
	userHandler     *userHandler
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
}

type UserRouteParams struct {
	fx.In

	R               *gin.Engine
	Logger          *logger.Logger
	UserHandler     *userHandler
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
}

func NewUserRoute(p UserRouteParams) *UserRoute {
	return &UserRoute{
		Engine:          p.R,
		logger:          p.Logger,
		userHandler:     p.UserHandler,
		authMiddleware:  p.AuthMiddleware,
		authzMiddleware: p.AuthzMiddleware,
	}
}

func (ur *UserRoute) Setup() {
	rg := ur.Group("/users")
	rg.POST("", ur.userHandler.createAccount)
	rg.POST("sign-in", ur.userHandler.signIn)
	rg.POST("/request-reset-password", ur.userHandler.requestResetPassword)
	rg.POST("/reset-password", ur.userHandler.resetPassword)

	rg.POST("/renew-token", ur.authMiddleware.RenewToken, ur.userHandler.renewToken)
	ur.registerAuthRoute(rg.Group("", ur.authMiddleware.RequireJwt))
}

func (ur *UserRoute) registerAuthRoute(rg *gin.RouterGroup) {
	requireUserUpdate := authzmiddleware.Require(ur.authzMiddleware, authorization.ActionUpdate, authorization.ObjectUser, authzmiddleware.UserParam("id"))

	rg.DELETE("/logout", ur.authMiddleware.InvalidateJwt, ur.userHandler.logout)
	rg.GET("/me", ur.userHandler.me)
	rg.PATCH("/:id", requireUserUpdate, ur.userHandler.updateUserProfile)
	rg.GET("/me/config", ur.userHandler.getMyConfig)
	rg.PATCH("/me/config", ur.userHandler.updateMyConfig)
	rg.GET("/me/followed-comics", ur.userHandler.getFollowedComics)

	requireRoleManage := authzmiddleware.Require(ur.authzMiddleware, authorization.ActionManage, authorization.ObjectRole)
	rg.GET("/:id/roles", requireRoleManage, ur.userHandler.getUserRoles)
	rg.POST("/:id/roles", requireRoleManage, ur.userHandler.assignUserRole)
	rg.DELETE("/:id/roles/:roleId", requireRoleManage, ur.userHandler.removeUserRole)
}
