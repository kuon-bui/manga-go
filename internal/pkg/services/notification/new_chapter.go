package notificationservice

import (
	"context"
	"errors"
	"fmt"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	notificationpkg "manga-go/internal/pkg/notification"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *NotificationService) handleComicNewChapter(ctx context.Context, payload *notificationpkg.FanoutPayload) error {
	chapter, err := s.chapterRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: payload.EntityID},
	}, nil)
	if err != nil {
		return err
	}

	if !chapter.IsPublished {
		return nil
	}

	comic, err := s.comicRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: chapter.ComicID},
	}, nil)
	if err != nil {
		return err
	}

	recipients, err := s.comicFollowRepo.FindFollowersByComicID(ctx, comic.ID)
	if err != nil {
		return err
	}

	eligibleRecipients := make([]*model.User, 0, len(recipients))
	sseEligible := make(map[uuid.UUID]bool)
	emailEligible := make(map[uuid.UUID]bool)
	for _, recipient := range recipients {
		if payload.TriggeredBy != nil && recipient.ID == *payload.TriggeredBy {
			continue
		}

		if !recipient.UserConfig.Has(model.UserConfigEnableComicNewChapterNotifications) {
			continue
		}

		eligibleRecipients = append(eligibleRecipients, recipient)
		sseEligible[recipient.ID] = recipient.UserConfig.Has(model.UserConfigEnableSSENotifications)
		emailEligible[recipient.ID] = recipient.UserConfig.Has(model.UserConfigEnableEmailNotifications)
	}

	if len(eligibleRecipients) == 0 {
		return nil
	}

	title := fmt.Sprintf("New chapter: %s", comic.Title)
	body := fmt.Sprintf("%s has a new chapter available: %s", comic.Title, s.chapterDisplayName(chapter))
	entityType := payload.EntityType
	dedupeKey := payload.DedupeKey
	notificationPayload := common.JSONMap{
		"comicId":       comic.ID,
		"comicSlug":     comic.Slug,
		"comicTitle":    comic.Title,
		"chapterId":     chapter.ID,
		"chapterSlug":   chapter.Slug,
		"chapterTitle":  chapter.Title,
		"chapterNumber": chapter.Number,
	}

	tx := s.gormDb.Begin().WithContext(ctx)
	if tx.Error != nil {
		return tx.Error
	}

	notificationRecord, err := s.notificationRepo.FindByDedupeKeyWithTransaction(tx, dedupeKey)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}

		notificationRecord = &model.Notification{
			Type:       payload.Type,
			Category:   notificationpkg.CategoryComic,
			ActorID:    payload.TriggeredBy,
			EntityType: &entityType,
			EntityID:   &payload.EntityID,
			DedupeKey:  &dedupeKey,
			Title:      title,
			Body:       body,
			Payload:    notificationPayload,
		}

		if err := s.notificationRepo.CreateWithTransaction(tx, notificationRecord); err != nil {
			tx.Rollback()
			return err
		}
	}

	userNotifications := make([]*model.UserNotification, 0, len(eligibleRecipients))
	userIDs := make([]uuid.UUID, 0, len(eligibleRecipients))
	for _, recipient := range eligibleRecipients {
		channelState := int64(0)
		if sseEligible[recipient.ID] {
			channelState |= notificationpkg.ChannelStateSSEQueued
		}
		if emailEligible[recipient.ID] {
			channelState |= notificationpkg.ChannelStateEmailQueued
		}

		userNotifications = append(userNotifications, &model.UserNotification{
			NotificationID: notificationRecord.ID,
			UserID:         recipient.ID,
			ChannelState:   channelState,
		})
		userIDs = append(userIDs, recipient.ID)
	}

	if err := s.userNotificationRepo.CreateListIgnoreConflictsWithTransaction(tx, userNotifications); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	persistedItems, err := s.userNotificationRepo.FindByNotificationAndUserIDs(ctx, notificationRecord.ID, userIDs)
	if err != nil {
		return err
	}

	if err := s.publishRealtime(ctx, persistedItems, sseEligible); err != nil {
		s.logger.Errorf("Failed to publish notification over SSE: %v", err)
	}

	if err := s.queueEmail(ctx, persistedItems, eligibleRecipients, emailEligible, comic, chapter); err != nil {
		s.logger.Errorf("Failed to queue notification email: %v", err)
	}

	return nil
}
