package model

import (
	"manga-go/internal/pkg/common"
	pknotification "manga-go/internal/pkg/notification"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

type Notification struct {
	common.SqlModel
	Type       pknotification.Type        `json:"type" gorm:"column:type"`
	Category   pknotification.Category    `json:"category" gorm:"column:category"`
	ActorID    *uuid.UUID                 `json:"actorId,omitempty" gorm:"column:actor_id"`
	EntityType *pknotification.EntityType `json:"entityType,omitempty" gorm:"column:entity_type"`
	EntityID   *uuid.UUID                 `json:"entityId,omitempty" gorm:"column:entity_id"`
	DedupeKey  *string                    `json:"-" gorm:"column:dedupe_key"`
	Title      string                     `json:"title" gorm:"column:title"`
	Body       string                     `json:"body" gorm:"column:body"`
	Payload    common.JSONMap             `json:"payload" gorm:"column:payload;type:jsonb"`

	Actor *User `json:"actor,omitempty" gorm:"foreignKey:ActorID"`
}

func (Notification) TableName() string {
	return "notifications"
}

func (n *Notification) Fake(f faker.Faker) {
	entityType := pknotification.EntityTypeChapter
	n.Type = pknotification.TypeComicNewChapter
	n.Category = pknotification.CategoryComic
	n.EntityType = &entityType
	n.Title = f.Lorem().Sentence(5)
	n.Body = f.Lorem().Paragraph(2)
	n.Payload = common.JSONMap{
		"headline": n.Title,
		"preview":  f.Lorem().Sentence(8),
	}
}
