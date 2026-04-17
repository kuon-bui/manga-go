package commentservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	commentrequest "manga-go/internal/pkg/request/comment"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *CommentService) CreateComment(ctx context.Context, userID uuid.UUID, req *commentrequest.CreateCommentRequest) response.Result {
	chapter, err := s.chapterRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: req.ChapterID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Chapter")
		}
		s.logger.Error("Failed to find chapter", "error", err)
		return response.ResultErrDb(err)
	}

	if req.ParentId != nil {
		_, err = s.commentRepo.FindOne(ctx, []any{
			clause.Eq{Column: "id", Value: *req.ParentId},
			clause.Eq{Column: "chapter_id", Value: chapter.ID},
			clause.Eq{Column: "comic_id", Value: chapter.ComicID},
		}, nil)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.ResultNotFound("Parent comment")
			}
			s.logger.Error("Failed to find parent comment", "error", err)
			return response.ResultErrDb(err)
		}
	}

	comment := model.Comment{
		UserId:    userID,
		ChapterId: chapter.ID,
		ComicId:   chapter.ComicID,
		ParentId:  req.ParentId,
		Content:   req.Content,
		PageIndex: req.PageIndex,
	}

	if err := s.commentRepo.Create(ctx, &comment); err != nil {
		s.logger.Error("Failed to create comment", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comment created successfully", comment)
}
