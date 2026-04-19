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
	// Validate: must have either comicId or chapterId
	if req.ComicID == nil && req.ChapterID == nil {
		return response.ResultError("Either comicId or chapterId is required")
	}

	var comment model.Comment
	comment.UserId = userID
	comment.ParentId = req.ParentId
	comment.Content = req.Content
	comment.PageIndex = req.PageIndex

	if req.ChapterID != nil {
		// Chapter-level or page-level comment
		chapter, err := s.chapterRepo.FindOne(ctx, []any{
			clause.Eq{Column: "id", Value: *req.ChapterID},
		}, nil)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.ResultNotFound("Chapter")
			}
			s.logger.Error("Failed to find chapter", "error", err)
			return response.ResultErrDb(err)
		}
		comment.ChapterId = &chapter.ID
		comment.ComicId = chapter.ComicID
	} else {
		// Comic-level comment: comicId only, no chapter
		comicID := *req.ComicID
		comment.ComicId = comicID
		comment.ChapterId = nil // no chapter for comic-level
	}

	if req.ParentId != nil {
		// Validate parent comment exists in same scope
		parentConditions := []any{
			clause.Eq{Column: "id", Value: *req.ParentId},
			clause.Eq{Column: "comic_id", Value: comment.ComicId},
		}
		_, err := s.commentRepo.FindOne(ctx, parentConditions, nil)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.ResultNotFound("Parent comment")
			}
			s.logger.Error("Failed to find parent comment", "error", err)
			return response.ResultErrDb(err)
		}
	}

	if err := s.commentRepo.Create(ctx, &comment); err != nil {
		s.logger.Error("Failed to create comment", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("Comment created successfully", comment)
}
