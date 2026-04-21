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
	"time"

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
	publishedAt := chapter.PublishedAt
	now := time.Now().UTC()

	err = s.chapterRepo.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		chapterUpdateData := map[string]any{
			"is_published": req.IsPublished,
		}

		switch {
		case req.IsPublished && !wasPublished:
			chapterUpdateData["published_at"] = now
		case !req.IsPublished && wasPublished:
			chapterUpdateData["published_at"] = nil
		}

		if err := s.chapterRepo.UpdateWithTransaction(tx, []any{
			clause.Eq{Column: "id", Value: chapter.ID},
		}, chapterUpdateData); err != nil {
			s.logger.Error("Failed to publish chapter", "error", err)
			return err
		}

		switch {
		case req.IsPublished && !wasPublished:
			if err := s.comicRepo.UpdateWithTransaction(tx, []any{
				clause.Eq{Column: "id", Value: comicID},
			}, map[string]any{
				"last_chapter_at": now,
			}); err != nil {
				s.logger.Error("Failed to update comic last_chapter_at", "error", err)
				return err
			}
		case !req.IsPublished && wasPublished:
			latestPublishedAt, err := s.chapterRepo.GetLatestPublishedChapterTimeByComicIDWithTransaction(tx, comicID)
			if err != nil {
				s.logger.Error("Failed to get latest published chapter time", "error", err)
				return err
			}

			if err := s.comicRepo.UpdateWithTransaction(tx, []any{
				clause.Eq{Column: "id", Value: comicID},
			}, map[string]any{
				"last_chapter_at": latestPublishedAt,
			}); err != nil {
				s.logger.Error("Failed to refresh comic last_chapter_at", "error", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		return response.ResultErrDb(err)
	}

	chapter.IsPublished = req.IsPublished
	switch {
	case req.IsPublished && !wasPublished:
		publishedAt = &now
	case !req.IsPublished && wasPublished:
		publishedAt = nil
	}
	chapter.PublishedAt = publishedAt

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
