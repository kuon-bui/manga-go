package commentroute

import (
	"manga-go/internal/pkg/logger"
	commentservice "manga-go/internal/pkg/services/comment"

	"go.uber.org/fx"
)

type CommentHandler struct {
	logger         *logger.Logger
	commentService *commentservice.CommentService
}

type CommentHandlerParams struct {
	fx.In
	Logger         *logger.Logger
	CommentService *commentservice.CommentService
}

func NewCommentHandler(p CommentHandlerParams) *CommentHandler {
	return &CommentHandler{
		logger:         p.Logger,
		commentService: p.CommentService,
	}
}
