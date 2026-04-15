package notification

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	CategoryComic = "comic"

	TypeComicNewChapter = "comic.new_chapter"

	EntityTypeChapter = "chapter"

	ChannelStateSSEQueued int64 = 1 << iota
	ChannelStateSSEDelivered
	ChannelStateEmailQueued
	ChannelStateEmailDelivered
	ChannelStateEmailFailed

	RedisUserChannelPrefix = "notifications:user"
)

type FanoutPayload struct {
	Type        string     `json:"type"`
	EntityType  string     `json:"entityType"`
	EntityID    uuid.UUID  `json:"entityId"`
	DedupeKey   string     `json:"dedupeKey"`
	TriggeredBy *uuid.UUID `json:"triggeredBy,omitempty"`
}

func UserChannel(userID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", RedisUserChannelPrefix, userID.String())
}
