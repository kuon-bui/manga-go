package commentservice

import (
	"manga-go/internal/pkg/logger"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	commentrepo "manga-go/internal/pkg/repo/comment"

	"go.uber.org/fx"
)

type CommentService struct {
	logger      *logger.Logger
	commentRepo *commentrepo.CommentRepo
	chapterRepo *chapterrepo.ChapterRepo
}

type CommentServiceParams struct {
	fx.In
	Logger      *logger.Logger
	CommentRepo *commentrepo.CommentRepo
	ChapterRepo *chapterrepo.ChapterRepo
}

func NewCommentService(p CommentServiceParams) *CommentService {
	return &CommentService{
		logger:      p.Logger,
		commentRepo: p.CommentRepo,
		chapterRepo: p.ChapterRepo,
	}
}
