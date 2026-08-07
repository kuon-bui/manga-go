package authorizationroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Route struct {
	engine          *gin.Engine
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
	handler         *Handler
}

type RouteParams struct {
	fx.In

	Engine          *gin.Engine
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
	Handler         *Handler
}

func NewRoute(params RouteParams) *Route {
	return &Route{
		engine:          params.Engine,
		authMiddleware:  params.AuthMiddleware,
		authzMiddleware: params.AuthzMiddleware,
		handler:         params.Handler,
	}
}

func (r *Route) Setup() {
	requireAuditRead := authzmiddleware.Require(
		r.authzMiddleware,
		authorization.ActionRead,
		authorization.ObjectAuditLog,
	)
	rg := r.engine.Group("/authorization", r.authMiddleware.RequireJwt)
	rg.GET("/audit-logs", requireAuditRead, r.handler.getAuditLogs)
}
