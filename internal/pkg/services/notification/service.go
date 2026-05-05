package notificationservice

import (
	"context"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/mail"
	notificationpkg "manga-go/internal/pkg/notification"
	pknotification "manga-go/internal/pkg/notification"
	"manga-go/internal/pkg/redis"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	comicfollowrepo "manga-go/internal/pkg/repo/comic_follow"
	notificationrepo "manga-go/internal/pkg/repo/notification"
	userrepo "manga-go/internal/pkg/repo/user"
	usernotificationrepo "manga-go/internal/pkg/repo/user_notification"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type NotificationItem struct {
	ID        string                  `json:"id"`
	Type      pknotification.Type     `json:"type"`
	Category  pknotification.Category `json:"category"`
	Title     string                  `json:"title"`
	Body      string                  `json:"body"`
	IsSeen    bool                    `json:"isSeen"`
	IsRead    bool                    `json:"isRead"`
	CreatedAt any                     `json:"createdAt"`
	Payload   any                     `json:"payload"`
}

type NotificationService struct {
	logger               *logger.Logger
	gormDb               *gorm.DB
	rds                  *redis.Redis
	asynqClient          *asynq.Client
	mailDialer           *mail.MailDialer
	notificationRepo     *notificationrepo.NotificationRepo
	userNotificationRepo *usernotificationrepo.UserNotificationRepo
	comicFollowRepo      *comicfollowrepo.ComicFollowRepo
	chapterRepo          *chapterrepo.ChapterRepo
	comicRepo            *comicrepo.ComicRepo
	userRepo             *userrepo.UserRepository
}

type NotificationServiceParams struct {
	fx.In
	Logger               *logger.Logger
	GormDb               *gorm.DB
	Redis                *redis.Redis
	AsynqClient          *asynq.Client
	MailDialer           *mail.MailDialer
	NotificationRepo     *notificationrepo.NotificationRepo
	UserNotificationRepo *usernotificationrepo.UserNotificationRepo
	ComicFollowRepo      *comicfollowrepo.ComicFollowRepo
	ChapterRepo          *chapterrepo.ChapterRepo
	ComicRepo            *comicrepo.ComicRepo
	UserRepo             *userrepo.UserRepository
}

func NewNotificationService(p NotificationServiceParams) *NotificationService {
	return &NotificationService{
		logger:               p.Logger,
		gormDb:               p.GormDb,
		rds:                  p.Redis,
		asynqClient:          p.AsynqClient,
		mailDialer:           p.MailDialer,
		notificationRepo:     p.NotificationRepo,
		userNotificationRepo: p.UserNotificationRepo,
		comicFollowRepo:      p.ComicFollowRepo,
		chapterRepo:          p.ChapterRepo,
		comicRepo:            p.ComicRepo,
		userRepo:             p.UserRepo,
	}
}

func (s *NotificationService) SubscribeUserChannel(ctx context.Context, userID uuid.UUID) *goredis.PubSub {
	return s.rds.Client().Subscribe(ctx, notificationpkg.UserChannel(userID))
}
