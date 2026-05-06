package commentservice

import (
	"context"
	"encoding/json"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	notificationpkg "manga-go/internal/pkg/notification"
	commentrequest "manga-go/internal/pkg/request/comment"
	queueconstant "manga-go/internal/queue/queue_constant"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
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

	if req.ParentId != nil && s.asynqClient != nil {
		payload, err := json.Marshal(notificationpkg.FanoutPayload{
			Type:        notificationpkg.TypeCommentNew,
			EntityType:  notificationpkg.EntityTypeComment,
			EntityID:    comment.ID,
			DedupeKey:   "comment-reply:" + comment.ID.String(),
			TriggeredBy: &userID,
		})
		if err != nil {
			s.logger.Error("Failed to marshal comment reply notification fanout payload", "error", err)
		} else {
			task := asynq.NewTask(queueconstant.NOTIFICATION_FANOUT_TASK, payload, asynq.MaxRetry(5))
			if _, err := s.asynqClient.Enqueue(task, asynq.Queue(queueconstant.NOTIFICATION_QUEUE)); err != nil {
				s.logger.Error("Failed to enqueue comment reply notification fanout task", "error", err)
			}
		}
	}

	return response.ResultSuccess("Comment created successfully", comment)
}
