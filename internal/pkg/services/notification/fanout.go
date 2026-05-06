package notificationservice

import (
	"context"
	"encoding/json"
	"fmt"
	"manga-go/internal/pkg/mail/mailable"
	"manga-go/internal/pkg/model"
	notificationpkg "manga-go/internal/pkg/notification"
	"time"

	"github.com/google/uuid"
)

func (s *NotificationService) HandleFanout(ctx context.Context, payload *notificationpkg.FanoutPayload) error {
	switch payload.Type {
	case notificationpkg.TypeComicNewChapter:
		return s.handleComicNewChapter(ctx, payload)
	case notificationpkg.TypeCommentNew:
		return s.handleCommentReply(ctx, payload)
	default:
		return fmt.Errorf("unsupported notification type: %s", payload.Type)
	}
}

func (s *NotificationService) publishRealtime(ctx context.Context, items []*model.UserNotification, sseEligible map[uuid.UUID]bool) error {
	var firstErr error
	for _, item := range items {
		if !sseEligible[item.UserID] {
			continue
		}

		payload, err := json.Marshal(s.mapNotificationItem(item))
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		if err := s.rds.Client().Publish(ctx, notificationpkg.UserChannel(item.UserID), payload).Err(); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		newState := item.ChannelState | notificationpkg.ChannelStateSSEDelivered
		if err := s.userNotificationRepo.MarkSSEDeliveredByID(ctx, item.ID, newState); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

func (s *NotificationService) queueEmail(ctx context.Context, items []*model.UserNotification, recipients []*model.User, emailEligible map[uuid.UUID]bool, comic *model.Comic, chapter *model.Chapter) error {
	itemByUserID := make(map[uuid.UUID]*model.UserNotification, len(items))
	for _, item := range items {
		itemByUserID[item.UserID] = item
	}

	outgoingMails := make([]mailable.MailableInterface, 0, len(recipients))
	queuedIDs := make([]uuid.UUID, 0, len(recipients))
	chapterDisplayName := s.chapterDisplayName(chapter)
	comicThumbnailURL := ""
	if comic.Thumbnail != nil {
		comicThumbnailURL = *comic.Thumbnail
	}
	for _, recipient := range recipients {
		if !emailEligible[recipient.ID] || recipient.Email == "" {
			continue
		}

		item, ok := itemByUserID[recipient.ID]
		if !ok || item.EmailedAt != nil || item.Notification == nil {
			continue
		}

		userName := recipient.Name
		if userName == "" {
			userName = recipient.Email
		}

		outgoingMail := mailable.NewComicNewChapterMail(mailable.ComicNewChapterMailParams{
			UserName:           userName,
			ComicTitle:         comic.Title,
			ComicThumbnailURL:  comicThumbnailURL,
			ChapterDisplayName: chapterDisplayName,
			ChapterNumber:      chapter.Number,
			ChapterTitle:       chapter.Title,
			CurrentYear:        time.Now().Year(),
		}).AddToFormat(mailable.MailAddress{
			Name:    recipient.Name,
			Address: recipient.Email,
		})

		outgoingMails = append(outgoingMails, outgoingMail)
		queuedIDs = append(queuedIDs, item.ID)
	}

	if len(outgoingMails) == 0 {
		return nil
	}

	if err := s.mailDialer.Dispatch(s.asynqClient, outgoingMails...); err != nil {
		return err
	}

	return s.userNotificationRepo.MarkEmailQueuedByIDs(ctx, queuedIDs)
}

func (s *NotificationService) chapterDisplayName(chapter *model.Chapter) string {
	if chapter.Title != "" {
		return chapter.Title
	}

	if chapter.Number != "" {
		return "Chapter " + chapter.Number
	}

	return chapter.Slug
}

func (s *NotificationService) notificationActorName(user *model.User) string {
	if user == nil {
		return ""
	}

	if user.Name != "" {
		return user.Name
	}

	return user.Email
}
