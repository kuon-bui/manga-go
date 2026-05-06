package notificationservice

import (
	"context"
	"fmt"
	"testing"
	"time"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	notificationpkg "manga-go/internal/pkg/notification"
	commentrepo "manga-go/internal/pkg/repo/comment"
	notificationrepo "manga-go/internal/pkg/repo/notification"
	userrepo "manga-go/internal/pkg/repo/user"
	usernotificationrepo "manga-go/internal/pkg/repo/user_notification"
	"manga-go/internal/pkg/testutil"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationTestSchema struct {
	testutil.SQLModel
	Type       notificationpkg.Type        `gorm:"column:type"`
	Category   notificationpkg.Category    `gorm:"column:category"`
	ActorID    *uuid.UUID                  `gorm:"column:actor_id;type:uuid"`
	EntityType *notificationpkg.EntityType `gorm:"column:entity_type"`
	EntityID   *uuid.UUID                  `gorm:"column:entity_id;type:uuid"`
	DedupeKey  *string                     `gorm:"column:dedupe_key"`
	Title      string                      `gorm:"column:title"`
	Body       string                      `gorm:"column:body"`
	Payload    common.JSONMap              `gorm:"column:payload"`
}

func (notificationTestSchema) TableName() string {
	return "notifications"
}

type userNotificationTestSchema struct {
	testutil.SQLModel
	NotificationID uuid.UUID  `gorm:"column:notification_id;type:uuid"`
	UserID         uuid.UUID  `gorm:"column:user_id;type:uuid"`
	ChannelState   int64      `gorm:"column:channel_state"`
	IsSeen         bool       `gorm:"column:is_seen"`
	SeenAt         *time.Time `gorm:"column:seen_at"`
	IsRead         bool       `gorm:"column:is_read"`
	ReadAt         *time.Time `gorm:"column:read_at"`
	EmailedAt      *time.Time `gorm:"column:emailed_at"`
	PushedAt       *time.Time `gorm:"column:pushed_at"`
}

func (userNotificationTestSchema) TableName() string {
	return "user_notifications"
}

func newNotificationServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.NewString())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	return db
}

func newNotificationServiceForFanoutTest(t *testing.T) (*NotificationService, *gorm.DB) {
	t.Helper()

	db := newNotificationServiceTestDB(t)
	testutil.MustSyncSchemas(t, db,
		&testutil.User{},
		&testutil.Comment{},
		&notificationTestSchema{},
		&userNotificationTestSchema{},
	)

	service := &NotificationService{
		logger:               logger.NewLogger(),
		gormDb:               db,
		notificationRepo:     notificationrepo.NewNotificationRepo(db),
		userNotificationRepo: usernotificationrepo.NewUserNotificationRepo(db),
		userRepo:             userrepo.NewUserRepository(db, nil),
		commentRepo:          commentrepo.NewCommentRepo(db),
	}

	return service, db
}

func TestHandleFanoutCommentReplyCreatesNotificationForParentAuthor(t *testing.T) {
	s, db := newNotificationServiceForFanoutTest(t)
	ctx := context.Background()

	recipientConfig := model.NewUserConfig()
	recipientConfig.Set(model.UserConfigEnableCommentReplyNotifications, true)
	recipientConfig.Set(model.UserConfigEnableEmailNotifications, true)

	recipientID := uuid.New()
	replierID := uuid.New()
	comicID := uuid.New()

	recipient := &testutil.User{
		SQLModel:   testutil.SQLModel{ID: recipientID},
		Name:       "Owner",
		Email:      "owner@example.com",
		UserConfig: []byte(recipientConfig),
	}
	replier := &testutil.User{
		SQLModel:   testutil.SQLModel{ID: replierID},
		Name:       "Alice",
		Email:      "alice@example.com",
		UserConfig: []byte(model.NewUserConfig()),
	}
	if err := db.Create(recipient).Error; err != nil {
		t.Fatalf("failed to create recipient: %v", err)
	}
	if err := db.Create(replier).Error; err != nil {
		t.Fatalf("failed to create replier: %v", err)
	}

	parentComment := &testutil.Comment{
		SQLModel: testutil.SQLModel{ID: uuid.New()},
		UserID:   recipientID,
		ComicID:  comicID,
		Content:  "parent comment",
	}
	if err := db.Create(parentComment).Error; err != nil {
		t.Fatalf("failed to create parent comment: %v", err)
	}

	replyComment := &testutil.Comment{
		SQLModel: testutil.SQLModel{ID: uuid.New()},
		UserID:   replierID,
		ComicID:  comicID,
		ParentID: &parentComment.ID,
		Content:  "reply comment",
	}
	if err := db.Create(replyComment).Error; err != nil {
		t.Fatalf("failed to create reply comment: %v", err)
	}

	if err := s.HandleFanout(ctx, &notificationpkg.FanoutPayload{
		Type:        notificationpkg.TypeCommentNew,
		EntityType:  notificationpkg.EntityTypeComment,
		EntityID:    replyComment.ID,
		DedupeKey:   "comment-reply:" + replyComment.ID.String(),
		TriggeredBy: &replierID,
	}); err != nil {
		t.Fatalf("expected comment reply fanout to succeed: %v", err)
	}

	var items []model.UserNotification
	if err := db.Preload("Notification").Find(&items).Error; err != nil {
		t.Fatalf("failed to load user notifications: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 user notification, got %d", len(items))
	}

	item := items[0]
	if item.UserID != recipientID {
		t.Fatalf("expected recipient %s, got %s", recipientID, item.UserID)
	}
	if item.ChannelState != 0 {
		t.Fatalf("expected comment reply notifications to avoid queued email delivery, got channel state %d", item.ChannelState)
	}
	if item.EmailedAt != nil {
		t.Fatal("expected comment reply notifications to skip email delivery")
	}
	if item.Notification == nil {
		t.Fatal("expected preloaded notification")
	}
	if item.Notification.Type != notificationpkg.TypeCommentNew {
		t.Fatalf("expected notification type %s, got %s", notificationpkg.TypeCommentNew, item.Notification.Type)
	}
	if item.Notification.Category != notificationpkg.CategoryComment {
		t.Fatalf("expected notification category %s, got %s", notificationpkg.CategoryComment, item.Notification.Category)
	}
	if item.Notification.EntityType == nil || *item.Notification.EntityType != notificationpkg.EntityTypeComment {
		t.Fatalf("expected entity type %s, got %#v", notificationpkg.EntityTypeComment, item.Notification.EntityType)
	}
	if item.Notification.EntityID == nil || *item.Notification.EntityID != replyComment.ID {
		t.Fatalf("expected entity id %s, got %#v", replyComment.ID, item.Notification.EntityID)
	}
	if item.Notification.Title != "New reply to your comment" {
		t.Fatalf("unexpected title: %s", item.Notification.Title)
	}
	if item.Notification.Body != "Alice replied to your comment" {
		t.Fatalf("unexpected body: %s", item.Notification.Body)
	}
	if fmt.Sprint(item.Notification.Payload["parentCommentId"]) != parentComment.ID.String() {
		t.Fatalf("expected payload parentCommentId %s, got %v", parentComment.ID, item.Notification.Payload["parentCommentId"])
	}
}

func TestHandleFanoutCommentReplySkipsSelfReplies(t *testing.T) {
	s, db := newNotificationServiceForFanoutTest(t)
	ctx := context.Background()

	config := model.NewUserConfig()
	config.Set(model.UserConfigEnableCommentReplyNotifications, true)

	userID := uuid.New()
	comicID := uuid.New()
	user := &testutil.User{
		SQLModel:   testutil.SQLModel{ID: userID},
		Name:       "Author",
		Email:      "author@example.com",
		UserConfig: []byte(config),
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	parentComment := &testutil.Comment{
		SQLModel: testutil.SQLModel{ID: uuid.New()},
		UserID:   userID,
		ComicID:  comicID,
		Content:  "parent comment",
	}
	if err := db.Create(parentComment).Error; err != nil {
		t.Fatalf("failed to create parent comment: %v", err)
	}

	replyComment := &testutil.Comment{
		SQLModel: testutil.SQLModel{ID: uuid.New()},
		UserID:   userID,
		ComicID:  comicID,
		ParentID: &parentComment.ID,
		Content:  "self reply",
	}
	if err := db.Create(replyComment).Error; err != nil {
		t.Fatalf("failed to create reply comment: %v", err)
	}

	if err := s.HandleFanout(ctx, &notificationpkg.FanoutPayload{
		Type:        notificationpkg.TypeCommentNew,
		EntityType:  notificationpkg.EntityTypeComment,
		EntityID:    replyComment.ID,
		DedupeKey:   "comment-reply:" + replyComment.ID.String(),
		TriggeredBy: &userID,
	}); err != nil {
		t.Fatalf("expected self reply fanout to be ignored without error: %v", err)
	}

	var total int64
	if err := db.Model(&model.UserNotification{}).Count(&total).Error; err != nil {
		t.Fatalf("failed to count user notifications: %v", err)
	}
	if total != 0 {
		t.Fatalf("expected no user notifications for self replies, got %d", total)
	}
}
