package commentroute

import (
	authmiddleware "manga-go/internal/app/middleware/auth"
	authzmiddleware "manga-go/internal/app/middleware/authz"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type CommentRoute struct {
	*gin.Engine
	logger          *logger.Logger
	authMiddleware  *authmiddleware.AuthMiddleware
	authzMiddleware *authzmiddleware.AuthzMiddleware
	commentHandler  *CommentHandler
}

type CommentRouteParams struct {
	fx.In
	R               *gin.Engine
	Logger          *logger.Logger
	AuthMiddleware  *authmiddleware.AuthMiddleware
	AuthzMiddleware *authzmiddleware.AuthzMiddleware
	CommentHandler  *CommentHandler
}

func NewCommentRoute(p CommentRouteParams) *CommentRoute {
	return &CommentRoute{
		logger:          p.Logger,
		Engine:          p.R,
		authMiddleware:  p.AuthMiddleware,
		authzMiddleware: p.AuthzMiddleware,
		commentHandler:  p.CommentHandler,
	}
}

func (cr *CommentRoute) Setup() {
	rg := cr.Group("/comments", cr.authMiddleware.RequireJwt)
	requireCommentCreate := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionCreate, authorization.ObjectComment)
	requireCommentUpdate := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionUpdate, authorization.ObjectComment, cr.authzMiddleware.CommentParam("id"))
	requireCommentDelete := authzmiddleware.Require(cr.authzMiddleware, authorization.ActionDelete, authorization.ObjectComment, cr.authzMiddleware.CommentParam("id"))

	rg.GET("", cr.commentHandler.getComments)
	rg.GET("/new", cr.commentHandler.getNewComments)
	rg.POST("", requireCommentCreate, cr.commentHandler.createComment)

	idRg := rg.Group("/:id")
	idRg.GET("", cr.commentHandler.getComment)
	idRg.GET("/replies", cr.commentHandler.getCommentReplies)
	idRg.PUT("", requireCommentUpdate, cr.commentHandler.updateComment)
	idRg.DELETE("", requireCommentDelete, cr.commentHandler.deleteComment)
	idRg.POST("/reactions", cr.commentHandler.handleReaction)
	idRg.POST("/report", cr.commentHandler.reportComment)
}
