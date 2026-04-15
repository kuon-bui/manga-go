package chapterserivce

import (
	"context"
	"encoding/json"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	notificationpkg "manga-go/internal/pkg/notification"
	chapterrequest "manga-go/internal/pkg/request/chapter"
	"manga-go/internal/pkg/utils"
	queueconstant "manga-go/internal/queue/queue_constant"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ChapterService) PublishChapter(ctx context.Context, chapterSlug string, req *chapterrequest.PublishChapterRequest) response.Result {
	comicID, ok := common.GetComicIdFromContext(ctx)
	if !ok {
		return response.ResultError("Comic not found in context")
	}

	chapter, err := s.chapterRepo.FindOne(ctx, []any{
		clause.Eq{Column: "comic_id", Value: comicID},
		clause.Eq{Column: "slug", Value: chapterSlug},
	}, map[string]common.MoreKeyOption{
		"Pages": {},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("Chapter")
		}
		s.logger.Error("Failed to find chapter", "error", err)
		return response.ResultErrDb(err)
	}

	if req.IsPublished && len(chapter.Pages) == 0 {
		return response.ResultError("Chapter must have at least one page before publishing")
	}

	wasPublished := chapter.IsPublished

	if err := s.chapterRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: chapter.ID},
	}, map[string]any{
		"is_published": req.IsPublished,
	}); err != nil {
		s.logger.Error("Failed to publish chapter", "error", err)
		return response.ResultErrDb(err)
	}

	chapter.IsPublished = req.IsPublished

	msg := "Chapter unpublished successfully"
	if req.IsPublished {
		msg = "Chapter published successfully"
		if !wasPublished {
			var triggeredBy *uuid.UUID
			if currentUser, err := utils.GetCurrentUserFormContext(ctx); err == nil {
				triggeredBy = &currentUser.ID
			}

			payload, err := json.Marshal(notificationpkg.FanoutPayload{
				Type:        notificationpkg.TypeComicNewChapter,
				EntityType:  notificationpkg.EntityTypeChapter,
				EntityID:    chapter.ID,
				DedupeKey:   "chapter-published:" + chapter.ID.String(),
				TriggeredBy: triggeredBy,
			})
			if err != nil {
				s.logger.Error("Failed to marshal notification fanout payload", "error", err)
				return response.ResultErrInternal(err)
			}

			task := asynq.NewTask(queueconstant.NOTIFICATION_FANOUT_TASK, payload, asynq.MaxRetry(5))
			if _, err := s.asynqClient.Enqueue(task, asynq.Queue(queueconstant.NOTIFICATION_QUEUE)); err != nil {
				s.logger.Error("Failed to enqueue notification fanout task", "error", err)
				return response.ResultErrInternal(err)
			}
		}
	}

	return response.ResultSuccess(msg, chapter)
}
