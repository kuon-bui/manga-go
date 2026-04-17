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
	rg.POST("", cr.commentHandler.createComment)

	idRg := rg.Group("/:id")
	idRg.GET("", cr.commentHandler.getComment)
	idRg.GET("/replies", cr.commentHandler.getCommentReplies)
	idRg.PUT("", cr.commentHandler.updateComment)
	idRg.DELETE("", cr.commentHandler.deleteComment)
	idRg.POST("/reactions", cr.commentHandler.handleReaction)
}
