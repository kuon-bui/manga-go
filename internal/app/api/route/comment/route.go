package commentroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type CommentRoute struct {
	logger         *logger.Logger
	r              *gin.Engine
	authMiddleware *authmiddleware.AuthMiddleware
	commentHandler *CommentHandler
}

type CommentRouteParams struct {
	fx.In
	Logger         *logger.Logger
	R              *gin.Engine
	AuthMiddleware *authmiddleware.AuthMiddleware
	CommentHandler *CommentHandler
}

func NewCommentRoute(p CommentRouteParams) *CommentRoute {
	return &CommentRoute{
		logger:         p.Logger,
		r:              p.R,
		authMiddleware: p.AuthMiddleware,
		commentHandler: p.CommentHandler,
	}
}

func (cr *CommentRoute) Setup() {
	rg := cr.r.Group("/comments", cr.authMiddleware.RequireJwt)

	rg.GET("", cr.commentHandler.getComments)
	rg.GET("/:id", cr.commentHandler.getComment)
	rg.POST("", cr.commentHandler.createComment)
	rg.PUT("/:id", cr.commentHandler.updateComment)
	rg.DELETE("/:id", cr.commentHandler.deleteComment)
}
