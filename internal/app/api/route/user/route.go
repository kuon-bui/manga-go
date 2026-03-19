package userroute

import (
	authmiddleware "base-go/internal/app/middleware/auth"
	"base-go/internal/pkg/logger"

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
	rg.POST("/", ur.userHandler.createAccount)
	rg.POST("sign-in", ur.userHandler.signIn)
	rg.POST("/request-reset-password", ur.userHandler.requestResetPassword)
	rg.POST("/reset-password", ur.userHandler.resetPassword)

	rg.POST("/renew-token", ur.authMiddleware.RenewToken, ur.userHandler.renewToken)
	ur.registerAuthRoute(rg.Group("", ur.authMiddleware.RequireJwt))
}

func (ur *UserRoute) registerAuthRoute(rg *gin.RouterGroup) {
	rg.DELETE("/logout", ur.authMiddleware.InvalidateJwt, ur.userHandler.logout)
	rg.GET("/me", ur.userHandler.me)

}
