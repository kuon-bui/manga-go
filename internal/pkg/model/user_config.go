package model

import (
	"database/sql/driver"
	"fmt"
)

const (
	UserConfigSeenNotificationCenter = iota
	UserConfigEnableSSENotifications
	UserConfigEnableEmailNotifications
	UserConfigEnableComicNewChapterNotifications
	UserConfigEnableCommentReplyNotifications
	UserConfigEnableMentionNotifications
	UserConfigEnableSystemAnnouncements
)

type UserConfig []byte

type UserConfigResponse struct {
	SeenNotificationCenter             bool `json:"seenNotificationCenter"`
	EnableSSENotifications             bool `json:"enableSseNotifications"`
	EnableEmailNotifications           bool `json:"enableEmailNotifications"`
	EnableComicNewChapterNotifications bool `json:"enableComicNewChapterNotifications"`
	EnableCommentReplyNotifications    bool `json:"enableCommentReplyNotifications"`
	EnableMentionNotifications         bool `json:"enableMentionNotifications"`
	EnableSystemAnnouncements          bool `json:"enableSystemAnnouncements"`
}

func NewUserConfig() UserConfig {
	return UserConfig{0}
}

func DefaultUserConfig() UserConfig {
	config := NewUserConfig()
	config.Set(UserConfigEnableSSENotifications, true)
	config.Set(UserConfigEnableComicNewChapterNotifications, true)

	return config
}

func (c UserConfig) Has(bit int) bool {
	bytePos := bit / 8
	bitPos := uint(bit % 8)

	if bytePos >= len(c) {
		return false
	}

	return c[bytePos]&(1<<bitPos) != 0
}

func (c *UserConfig) Set(bit int, enabled bool) {
	c.ensureSize(bit)

	bytePos := bit / 8
	bitPos := uint(bit % 8)

	if enabled {
		(*c)[bytePos] |= 1 << bitPos
		return
	}

	(*c)[bytePos] &^= 1 << bitPos
}

func (c *UserConfig) ensureSize(bit int) {
	bytePos := bit / 8
	for len(*c) <= bytePos {
		*c = append(*c, 0)
	}

	if len(*c) == 0 {
		*c = append(*c, 0)
	}
}

func (c *UserConfig) Scan(value any) error {
	if value == nil {
		*c = NewUserConfig()
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("UserConfig.Scan: expected []byte, got %T", value)
	}

	if len(bytes) == 0 {
		*c = NewUserConfig()
		return nil
	}

	cp := make([]byte, len(bytes))
	copy(cp, bytes)
	*c = UserConfig(cp)

	return nil
}

func (c UserConfig) Value() (driver.Value, error) {
	if len(c) == 0 {
		return []byte{0}, nil
	}

	return []byte(c), nil
}

func (c UserConfig) ToResponse() UserConfigResponse {
	return UserConfigResponse{
		SeenNotificationCenter:             c.Has(UserConfigSeenNotificationCenter),
		EnableSSENotifications:             c.Has(UserConfigEnableSSENotifications),
		EnableEmailNotifications:           c.Has(UserConfigEnableEmailNotifications),
		EnableComicNewChapterNotifications: c.Has(UserConfigEnableComicNewChapterNotifications),
		EnableCommentReplyNotifications:    c.Has(UserConfigEnableCommentReplyNotifications),
		EnableMentionNotifications:         c.Has(UserConfigEnableMentionNotifications),
		EnableSystemAnnouncements:          c.Has(UserConfigEnableSystemAnnouncements),
	}
}
