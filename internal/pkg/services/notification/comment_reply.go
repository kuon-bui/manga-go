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

func (s *NotificationService) handleCommentReply(ctx context.Context, payload *notificationpkg.FanoutPayload) error {
	fmt.Println("Handling comment reply notification fanout for entity ID:", payload.EntityID)
	replyComment, err := s.commentRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: payload.EntityID},
	}, map[string]common.MoreKeyOption{
		"User": {},
	})
	if err != nil {
		return err
	}

	if replyComment.ParentId == nil {
		return nil
	}

	parentComment, err := s.commentRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: *replyComment.ParentId},
	}, nil)
	if err != nil {
		return err
	}

	fmt.Println(payload.TriggeredBy, "is replying to comment by user", parentComment.UserId)
	fmt.Printf("parentComment.UserId=%s\n", parentComment.UserId)

	if payload.TriggeredBy != nil && parentComment.UserId == *payload.TriggeredBy {
		return nil
	}

	recipient, err := s.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: parentComment.UserId},
	}, nil)
	if err != nil {
		return err
	}

	if !recipient.UserConfig.Has(model.UserConfigEnableCommentReplyNotifications) {
		return nil
	}

	sseEligible := recipient.UserConfig.Has(model.UserConfigEnableSSENotifications)
	userIDs := []uuid.UUID{recipient.ID}
	title := "New reply to your comment"
	body := "Someone replied to your comment"
	if actorName := s.notificationActorName(replyComment.User); actorName != "" {
		body = fmt.Sprintf("%s replied to your comment", actorName)
	}

	entityType := payload.EntityType
	dedupeKey := payload.DedupeKey
	notificationPayload := common.JSONMap{
		"commentId":       replyComment.ID,
		"parentCommentId": parentComment.ID,
		"comicId":         replyComment.ComicId,
	}
	if replyComment.ChapterId != nil {
		notificationPayload["chapterId"] = *replyComment.ChapterId
	}
	if replyComment.PageIndex != nil {
		notificationPayload["pageIndex"] = *replyComment.PageIndex
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
			Category:   notificationpkg.CategoryComment,
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

	channelState := int64(0)
	if sseEligible {
		channelState |= notificationpkg.ChannelStateSSEQueued
	}

	userNotifications := []*model.UserNotification{{
		NotificationID: notificationRecord.ID,
		UserID:         recipient.ID,
		ChannelState:   channelState,
	}}

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

	if err := s.publishRealtime(ctx, persistedItems, map[uuid.UUID]bool{recipient.ID: sseEligible}); err != nil {
		s.logger.Errorf("Failed to publish comment reply notification over SSE: %v", err)
	}

	return nil
}
