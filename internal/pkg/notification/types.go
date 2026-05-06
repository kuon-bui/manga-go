package notification

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Category string

type Type string

type EntityType string

const (
	CategoryComic   Category = "comic"
	CategoryComment Category = "comment"
	CategorySystem  Category = "system"

	TypeComicNewChapter Type = "comic.new_chapter"
	TypeCommentNew      Type = "comment.reply"
	TypeSystemTest      Type = "system.test"

	EntityTypeChapter EntityType = "chapter"
	EntityTypeComment EntityType = "comment"

	ChannelStateSSEQueued int64 = 1 << iota
	ChannelStateSSEDelivered
	ChannelStateEmailQueued
	ChannelStateEmailDelivered
	ChannelStateEmailFailed

	RedisUserChannelPrefix = "notifications:user"
)

type FanoutPayload struct {
	Type        Type       `json:"type"`
	EntityType  EntityType `json:"entityType"`
	EntityID    uuid.UUID  `json:"entityId"`
	DedupeKey   string     `json:"dedupeKey"`
	TriggeredBy *uuid.UUID `json:"triggeredBy,omitempty"`
}

func UserChannel(userID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", RedisUserChannelPrefix, userID.String())
}

func GetAllCategories() []Category {
	return []Category{
		CategoryComic,
		CategoryComment,
		CategorySystem,
	}
}

func GetAllTypes() []Type {
	return []Type{
		TypeComicNewChapter,
		TypeCommentNew,
		TypeSystemTest,
	}
}

func GetAllEntityTypes() []EntityType {
	return []EntityType{
		EntityTypeChapter,
		EntityTypeComment,
	}
}

func CategoryValidationMessage(field string) string {
	allCategories := GetAllCategories()
	res := strings.Builder{}
	res.WriteString(string(allCategories[0]))
	for _, c := range allCategories[1:] {
		res.WriteString(", ")
		res.WriteString(string(c))
	}

	return fmt.Sprintf(
		"%s must be a valid notification category (%s)",
		field,
		res.String(),
	)
}

func TypeValidationMessage(field string) string {
	allTypes := GetAllTypes()
	res := strings.Builder{}
	res.WriteString(string(allTypes[0]))
	for _, t := range allTypes[1:] {
		res.WriteString(", ")
		res.WriteString(string(t))
	}

	return fmt.Sprintf(
		"%s must be a valid notification type (%s)",
		field,
		res.String(),
	)
}

func EntityTypeValidationMessage(field string) string {
	allEntityTypes := GetAllEntityTypes()
	res := strings.Builder{}
	res.WriteString(string(allEntityTypes[0]))
	for _, t := range allEntityTypes[1:] {
		res.WriteString(", ")
		res.WriteString(string(t))
	}

	return fmt.Sprintf(
		"%s must be a valid notification entity type (%s)",
		field,
		res.String(),
	)
}

func GetAllowedCategories() map[Category]bool {
	return map[Category]bool{
		CategoryComic:   true,
		CategoryComment: true,
		CategorySystem:  true,
	}
}

func GetAllowedTypes() map[Type]bool {
	return map[Type]bool{
		TypeComicNewChapter: true,
		TypeCommentNew:      true,
		TypeSystemTest:      true,
	}
}

func GetAllowedEntityTypes() map[EntityType]bool {
	return map[EntityType]bool{
		EntityTypeChapter: true,
		EntityTypeComment: true,
	}
}
