package userroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type UserRoute struct {
	logger         *logger.Logger
	r              *gin.Engine
	userHandler    *userHandler
	authMiddleware *authmiddleware.AuthMiddleware
}

type UserRouteParams struct {
	fx.In

	R              *gin.Engine
	Logger         *logger.Logger
	UserHandler    *userHandler
	AuthMiddleware *authmiddleware.AuthMiddleware
}

func NewUserRoute(p UserRouteParams) *UserRoute {
	return &UserRoute{
		r:              p.R,
		logger:         p.Logger,
		userHandler:    p.UserHandler,
		authMiddleware: p.AuthMiddleware,
	}
}

func (ur *UserRoute) Setup() {
	rg := ur.r.Group("/users")
	rg.POST("", ur.userHandler.createAccount)
	rg.POST("sign-in", ur.userHandler.signIn)
	rg.POST("/request-reset-password", ur.userHandler.requestResetPassword)
	rg.POST("/reset-password", ur.userHandler.resetPassword)

	rg.POST("/renew-token", ur.authMiddleware.RenewToken, ur.userHandler.renewToken)
	ur.registerAuthRoute(rg.Group("", ur.authMiddleware.RequireJwt))
}

func (ur *UserRoute) registerAuthRoute(rg *gin.RouterGroup) {
	rg.DELETE("/logout", ur.authMiddleware.InvalidateJwt, ur.userHandler.logout)
	rg.GET("/me", ur.userHandler.me)
	rg.GET("/me/config", ur.userHandler.getMyConfig)
	rg.PATCH("/me/config", ur.userHandler.updateMyConfig)
	rg.GET("/me/followed-comics", ur.userHandler.getFollowedComics)
	rg.GET("/:id/roles", ur.userHandler.getUserRoles)
	rg.POST("/:id/roles", ur.userHandler.assignUserRole)
	rg.DELETE("/:id/roles/:roleId", ur.userHandler.removeUserRole)
}
