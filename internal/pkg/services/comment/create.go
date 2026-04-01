package commentservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	commentrequest "manga-go/internal/pkg/request/comment"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (s *CommentService) CreateComment(ctx context.Context, userID uuid.UUID, req *commentrequest.CreateCommentRequest) response.Result {
	chapter, err := s.chapterRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: req.ChapterID},
	}, nil)
	if err != nil {
		s.logger.Error("Failed to find chapter", "error", err)
		return response.ResultErrDb(err)
	}

	comment := model.Comment{
		UserId:    userID,
		ChapterId: chapter.ID,
		ComicId:   chapter.ComicID,
		Content:   req.Content,
		PageIndex: req.PageIndex,
	}

	if err := s.commentRepo.Create(ctx, &comment); err != nil {
		s.logger.Error("Failed to create comment", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comment created successfully", comment)
}
