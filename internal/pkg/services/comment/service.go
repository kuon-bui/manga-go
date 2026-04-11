package commentservice

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	commentrepo "manga-go/internal/pkg/repo/comment"
	reactionrepo "manga-go/internal/pkg/repo/reaction"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

// CommentRepository defines the data access interface for Comment.
type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Comment, error)
	Update(ctx context.Context, conditions []any, data map[string]any) error
	DeleteSoft(ctx context.Context, conditions []any) error
	FindPaginated(ctx context.Context, conditions []any, paging *common.Paging, moreKeys map[string]common.MoreKeyOption) ([]*model.Comment, int64, error)
}

// ChapterRepository defines the subset of ChapterRepo used by CommentService.
type ChapterRepository interface {
	FindOne(ctx context.Context, conditions []any, moreKeys map[string]common.MoreKeyOption) (*model.Chapter, error)
}

// ReactionRepository defines the data access interface for Reaction.
type ReactionRepository interface {
	Create(ctx context.Context, reaction *model.Reaction) error
	DeleteSoft(ctx context.Context, conditions []any) error
	ExistsByCommentIdAndUserId(ctx context.Context, commentId, userId uuid.UUID) (bool, error)
}

type CommentService struct {
	logger       *logger.Logger
	commentRepo  CommentRepository
	chapterRepo  ChapterRepository
	reactionRepo ReactionRepository
}

type CommentServiceParams struct {
	fx.In
	Logger       *logger.Logger
	CommentRepo  *commentrepo.CommentRepo
	ChapterRepo  *chapterrepo.ChapterRepo
	ReactionRepo *reactionrepo.ReactionRepo
}

func NewCommentService(p CommentServiceParams) *CommentService {
	return &CommentService{
		logger:       p.Logger,
		commentRepo:  p.CommentRepo,
		chapterRepo:  p.ChapterRepo,
		reactionRepo: p.ReactionRepo,
	}
}

// NewCommentServiceWithRepos creates a CommentService with explicit repository interfaces,
// useful for unit testing.
func NewCommentServiceWithRepos(l *logger.Logger, commentRepo CommentRepository, chapterRepo ChapterRepository, reactionRepo ReactionRepository) *CommentService {
	return &CommentService{
		logger:       l,
		commentRepo:  commentRepo,
		chapterRepo:  chapterRepo,
		reactionRepo: reactionRepo,
	}
}
