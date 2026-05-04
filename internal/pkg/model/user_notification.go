package model

import (
	"manga-go/internal/pkg/common"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

type UserNotification struct {
	common.SqlModel
	NotificationID uuid.UUID  `json:"notificationId" gorm:"column:notification_id"`
	UserID         uuid.UUID  `json:"userId" gorm:"column:user_id"`
	ChannelState   int64      `json:"-" gorm:"column:channel_state"`
	IsSeen         bool       `json:"isSeen" gorm:"column:is_seen"`
	SeenAt         *time.Time `json:"seenAt,omitempty" gorm:"column:seen_at"`
	IsRead         bool       `json:"isRead" gorm:"column:is_read"`
	ReadAt         *time.Time `json:"readAt,omitempty" gorm:"column:read_at"`
	EmailedAt      *time.Time `json:"emailedAt,omitempty" gorm:"column:emailed_at"`
	PushedAt       *time.Time `json:"pushedAt,omitempty" gorm:"column:pushed_at"`

	Notification *Notification `json:"notification,omitempty" gorm:"foreignKey:NotificationID"`
	User         *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (UserNotification) TableName() string {
	return "user_notifications"
}

func (u *UserNotification) Fake(f faker.Faker) {
	u.IsSeen = f.Bool()
	u.IsRead = u.IsSeen && f.Bool()

	if u.IsSeen {
		seenAt := time.Now().Add(-time.Duration(f.IntBetween(1, 72)) * time.Hour)
		u.SeenAt = &seenAt
	}

	if u.IsRead {
		readAt := time.Now().Add(-time.Duration(f.IntBetween(1, 48)) * time.Hour)
		u.ReadAt = &readAt
	}
}
