package notificationseeder

import (
	"context"
	"errors"
	"fmt"
	"manga-go/internal/pkg/model"
	pknotification "manga-go/internal/pkg/notification"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	comicfollowrepo "manga-go/internal/pkg/repo/comic_follow"
	notificationrepo "manga-go/internal/pkg/repo/notification"
	userrepo "manga-go/internal/pkg/repo/user"
	usernotificationrepo "manga-go/internal/pkg/repo/user_notification"
	seederutil "manga-go/internal/pkg/seeder/util"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NotificationSeeder struct {
	notificationRepo     *notificationrepo.NotificationRepo
	userNotificationRepo *usernotificationrepo.UserNotificationRepo
	comicFollowRepo      *comicfollowrepo.ComicFollowRepo
	comicRepo            *comicrepo.ComicRepo
	chapterRepo          *chapterrepo.ChapterRepo
	userRepo             *userrepo.UserRepository
	faker                faker.Faker
}

func NewNotificationSeeder(
	notificationRepo *notificationrepo.NotificationRepo,
	userNotificationRepo *usernotificationrepo.UserNotificationRepo,
	comicFollowRepo *comicfollowrepo.ComicFollowRepo,
	comicRepo *comicrepo.ComicRepo,
	chapterRepo *chapterrepo.ChapterRepo,
	userRepo *userrepo.UserRepository,
	faker faker.Faker,
) *NotificationSeeder {
	return &NotificationSeeder{
		notificationRepo:     notificationRepo,
		userNotificationRepo: userNotificationRepo,
		comicFollowRepo:      comicFollowRepo,
		comicRepo:            comicRepo,
		chapterRepo:          chapterRepo,
		userRepo:             userRepo,
		faker:                faker,
	}
}

func (s *NotificationSeeder) Name() string {
	return "NotificationSeeder"
}

func (s *NotificationSeeder) Truncate(tx *gorm.DB) error {
	return seederutil.TruncateTables(tx, "user_notifications", "notifications")
}

func (s *NotificationSeeder) Seed(tx *gorm.DB) error {
	comics, err := s.comicRepo.FindAllWithTx(tx, []any{func(db *gorm.DB) *gorm.DB {
		return db.Order("title ASC")
	}}, nil)
	if err != nil {
		return err
	}
	users, err := s.userRepo.FindAllWithTx(tx, []any{func(db *gorm.DB) *gorm.DB {
		return db.Order("email ASC")
	}}, nil)
	if err != nil {
		return err
	}
	if len(comics) == 0 || len(users) == 0 {
		return nil
	}

	for index, comic := range comics {
		chapters, err := s.chapterRepo.FindAllWithTx(tx, []any{
			clause.Eq{Column: "comic_id", Value: comic.ID},
			func(db *gorm.DB) *gorm.DB { return db.Order("chapter_idx DESC") },
		}, nil)
		if err != nil {
			return err
		}
		if len(chapters) == 0 {
			continue
		}

		chapter := chapters[0]
		dedupeKey := fmt.Sprintf("seed:comic.new_chapter:%s", chapter.ID)
		notification, err := s.notificationRepo.FindByDedupeKeyWithTransaction(tx, dedupeKey)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notification = &model.Notification{}
			notification.Fake(s.faker)
			notification.Title = fmt.Sprintf("New chapter for %s", comic.Title)
			notification.Body = fmt.Sprintf("Chapter %s - %s is now available.", chapter.Number, chapter.Title)
			notification.DedupeKey = &dedupeKey
			notification.EntityID = &chapter.ID
			notification.Payload["comicId"] = comic.ID.String()
			notification.Payload["comicTitle"] = comic.Title
			notification.Payload["chapterId"] = chapter.ID.String()
			notification.Payload["chapterTitle"] = chapter.Title
			actorID := users[index%len(users)].ID
			notification.ActorID = &actorID
			if err := s.notificationRepo.CreateWithTransaction(tx, notification); err != nil {
				return err
			}
		}

		followers, err := s.comicFollowRepo.FindFollowersByComicID(context.Background(), comic.ID)
		if err != nil {
			return err
		}
		if len(followers) == 0 {
			followers = users[:min(3, len(users))]
		}

		items := make([]*model.UserNotification, 0, len(followers))
		for _, follower := range followers {
			item := &model.UserNotification{NotificationID: notification.ID, UserID: follower.ID, ChannelState: pknotification.ChannelStateSSEQueued}
			item.Fake(s.faker)
			items = append(items, item)
		}

		if err := s.userNotificationRepo.CreateListIgnoreConflictsWithTransaction(tx, items); err != nil {
			return err
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
