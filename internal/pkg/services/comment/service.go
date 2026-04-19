package commentservice

import (
	"manga-go/internal/pkg/logger"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	commentrepo "manga-go/internal/pkg/repo/comment"
	commentreportrepo "manga-go/internal/pkg/repo/comment_report"
	reactionrepo "manga-go/internal/pkg/repo/reaction"

	"go.uber.org/fx"
)

type CommentService struct {
	logger            *logger.Logger
	commentRepo       *commentrepo.CommentRepo
	chapterRepo       *chapterrepo.ChapterRepo
	reactionRepo      *reactionrepo.ReactionRepo
	commentReportRepo *commentreportrepo.CommentReportRepo
}

type CommentServiceParams struct {
	fx.In
	Logger            *logger.Logger
	CommentRepo       *commentrepo.CommentRepo
	ChapterRepo       *chapterrepo.ChapterRepo
	ReactionRepo      *reactionrepo.ReactionRepo
	CommentReportRepo *commentreportrepo.CommentReportRepo
}

func NewCommentService(p CommentServiceParams) *CommentService {
	return &CommentService{
		logger:            p.Logger,
		commentRepo:       p.CommentRepo,
		chapterRepo:       p.ChapterRepo,
		reactionRepo:      p.ReactionRepo,
		commentReportRepo: p.CommentReportRepo,
	}
}
