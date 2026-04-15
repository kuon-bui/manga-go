package model

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type Notification struct {
	common.SqlModel
	Type       string         `json:"type" gorm:"column:type"`
	Category   string         `json:"category" gorm:"column:category"`
	ActorID    *uuid.UUID     `json:"actorId,omitempty" gorm:"column:actor_id"`
	EntityType *string        `json:"entityType,omitempty" gorm:"column:entity_type"`
	EntityID   *uuid.UUID     `json:"entityId,omitempty" gorm:"column:entity_id"`
	DedupeKey  *string        `json:"-" gorm:"column:dedupe_key"`
	Title      string         `json:"title" gorm:"column:title"`
	Body       string         `json:"body" gorm:"column:body"`
	Payload    common.JSONMap `json:"payload" gorm:"column:payload;type:jsonb"`

	Actor *User `json:"actor,omitempty" gorm:"foreignKey:ActorID"`
}

func (Notification) TableName() string {
	return "notifications"
}
