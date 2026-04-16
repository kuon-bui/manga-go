# Notification System Technical Design

## 1. Overview

This document describes a generic notification architecture for manga-go that can support current and future notification types without redesigning the core model each time.

Primary use case for the first rollout:
- notify followers when a new chapter becomes published
- deliver in real time through SSE
- optionally deliver through email based on user settings

The design is intentionally generic so it can later support:
- comment reply notifications
- mention notifications
- system announcements
- moderation or admin notices
- role or permission related notifications

## 2. Goals

- Provide a reusable notification domain model.
- Keep notification history persisted in PostgreSQL.
- Support per-user state such as seen, read, and channel delivery state.
- Support real-time delivery through SSE.
- Support asynchronous email delivery through the existing Asynq + mail stack.
- Store durable user-level notification settings in the users table using a bitmap-backed custom type.
- Keep the publish chapter API fast by moving fan-out and email work to background jobs.

## 3. Non-Goals

- Mobile push notification delivery in the first phase.
- Cross-device per-session presence tracking in the first phase.
- Full preference matrix for every notification type in the first phase.
- Notification templating UI in the first phase.

## 4. Current State

Existing building blocks already present in the repository:

- Comic follow relationship is available through comic_follows.
- Chapter publish workflow is available in the chapter service.
- Redis is already available.
- Asynq client and worker infrastructure already exist.
- Generic mail dispatch already exists.
- Auth middleware already injects the current user into the context.

Missing pieces:

- generic notification event model
- per-user notification state model
- user notification settings bitmap in users
- SSE streaming endpoint and broker
- background fan-out pipeline for notification recipients

## 5. High-Level Architecture

The system is split into three layers:

1. Notification event
- Represents the source event only once.
- Example: a chapter was published.

2. User notification
- Represents the relationship between one event and one recipient user.
- Stores per-user state such as seen, read, and delivered channels.

3. User config bitmap
- Stored directly on users.
- Stores stable user-level preferences and lightweight UI flags.
- Does not store per-notification read state.

### 5.1 Architecture Flow

1. Chapter transitions from unpublished to published.
2. Chapter service enqueues a notification fan-out task.
3. Fan-out worker creates one notification event record if not already created.
4. Worker finds recipient users for that event.
5. Worker inserts user_notifications in batch.
6. Worker publishes SSE messages to Redis Pub/Sub channels per recipient.
7. Worker schedules email delivery for recipients whose config allows email.

## 6. Domain Model

### 6.1 Notification Event

The event record is generic and reusable.

Suggested fields:

- id: UUID
- type: string
- category: string
- actor_id: UUID nullable
- entity_type: string
- entity_id: UUID nullable
- dedupe_key: string nullable
- title: string
- body: text
- payload: jsonb
- created_at
- updated_at
- deleted_at

Semantics:

- type identifies the business event, for example comic.new_chapter
- category groups similar events, for example comic, comment, system
- actor_id is the user who caused the event if applicable
- entity_type and entity_id point to the main object related to the event
- dedupe_key is used for idempotency
- payload stores type-specific data used by clients and email templates

### 6.2 User Notification

This record stores recipient-specific state.

Suggested fields:

- id: UUID
- notification_id: UUID
- user_id: UUID
- channel_state: bigint default 0
- is_seen: bool default false
- seen_at: timestamptz nullable
- is_read: bool default false
- read_at: timestamptz nullable
- emailed_at: timestamptz nullable
- pushed_at: timestamptz nullable
- created_at
- updated_at
- deleted_at

Semantics:

- one row per user per notification
- is_seen is for inbox badge and first visual exposure
- is_read is for explicit user acknowledgement
- channel_state is a compact bitmap for delivery channel state if needed
- emailed_at is useful for auditing and retries

### 6.3 User Config Bitmap

User-level settings are stored in users.user_config.

Suggested storage:

- PostgreSQL type: bytea
- Go type: custom named type such as UserConfigFlags ReadBitset

Reasoning:

- user settings are low cardinality and stable
- bytea use to store byte data
- query and update operations are simpler than bytea for this use case
- the repository already uses custom DB-backed types, so adding one more custom type fits the codebase

This bitmap should only store:

- preference toggles
- feature flags
- lightweight UI state

This bitmap should not store:

- read state of individual notifications
- notification inbox history
- delivery logs

## 7. Suggested Notification Types

Use semantic string types instead of integer enums.

Initial naming convention:

- comic.new_chapter
- comic.updated
- comment.reply
- comment.mention
- user.role_assigned
- system.announcement

Rules:

- use lower case with dot-separated namespaces
- do not encode delivery channel in the type name
- keep type stable for analytics and routing logic

## 8. Suggested User Config Flag Map

Store these in users.user_config.

Suggested initial bit allocation:

- bit 0: notification center has been seen at least once
- bit 1: enable SSE notifications
- bit 2: enable email notifications
- bit 3: enable comic new chapter notifications
- bit 4: enable comment reply notifications
- bit 5: enable mention notifications
- bit 6: enable system announcements

Suggested Go API:

```go
type UserConfigFlags struct {
  SeenNotificationCenter             bool `json:"seenNotificationCenter"`
	EnableSSENotifications             bool `json:"enableSseNotifications"`
	EnableEmailNotifications           bool `json:"enableEmailNotifications"`
	EnableComicNewChapterNotifications bool `json:"enableComicNewChapterNotifications"`
	EnableCommentReplyNotifications    bool `json:"enableCommentReplyNotifications"`
	EnableMentionNotifications         bool `json:"enableMentionNotifications"`
	EnableSystemAnnouncements          bool `json:"enableSystemAnnouncements"`

}

const (
	UserConfigSeenNotificationCenter int = 1
	UserConfigEnableSSENotifications 
	UserConfigEnableEmailNotifications
	UserConfigEnableComicNewChapterNotifications
	UserConfigEnableCommentReplyNotifications
	UserConfigEnableMentionNotifications
	UserConfigEnableSystemAnnouncements
)
```

Suggested helper methods:
- read from bitmask data 
- store data as bitmask

Suggested response DTO:

```go
type UserConfigResponse struct {
	SeenNotificationCenter             bool `json:"seenNotificationCenter"`
	EnableSSENotifications             bool `json:"enableSseNotifications"`
	EnableEmailNotifications           bool `json:"enableEmailNotifications"`
	EnableComicNewChapterNotifications bool `json:"enableComicNewChapterNotifications"`
	EnableCommentReplyNotifications    bool `json:"enableCommentReplyNotifications"`
	EnableMentionNotifications         bool `json:"enableMentionNotifications"`
	EnableSystemAnnouncements          bool `json:"enableSystemAnnouncements"`
}
```

## 9. Database Design

### 9.1 users

Add column:

- user_config bytea not null default 0

Suggested migration:

```sql
ALTER TABLE users
ADD COLUMN user_config bytea NOT NULL DEFAULT 0;
```

### 9.2 notifications

Suggested table:

```sql
CREATE TABLE notifications (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    type VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    actor_id uuid NULL,
    entity_type VARCHAR(50) NULL,
    entity_id uuid NULL,
    dedupe_key VARCHAR(255) NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX ux_notifications_dedupe_key
ON notifications (dedupe_key)
WHERE dedupe_key IS NOT NULL AND deleted_at IS NULL;

CREATE INDEX idx_notifications_type_created_at
ON notifications (type, created_at DESC)
WHERE deleted_at IS NULL;
```

Notes:

- dedupe_key is optional because not every type needs it
- for comic.new_chapter, dedupe_key can be chapter-published:{chapter_id}

### 9.3 user_notifications

Suggested table:

```sql
CREATE TABLE user_notifications (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    notification_id uuid NOT NULL,
    user_id uuid NOT NULL,
    channel_state BIGINT NOT NULL DEFAULT 0,
    is_seen BOOLEAN NOT NULL DEFAULT FALSE,
    seen_at TIMESTAMPTZ NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMPTZ NULL,
    emailed_at TIMESTAMPTZ NULL,
    pushed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_user_notifications_notification
        FOREIGN KEY (notification_id) REFERENCES notifications(id),
    CONSTRAINT fk_user_notifications_user
        FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE UNIQUE INDEX ux_user_notifications_notification_user
ON user_notifications (notification_id, user_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_user_notifications_user_created_at
ON user_notifications (user_id, created_at DESC)
WHERE deleted_at IS NULL;

CREATE INDEX idx_user_notifications_user_unread
ON user_notifications (user_id, is_read, is_seen, created_at DESC)
WHERE deleted_at IS NULL;
```

## 10. Channel State Bitmap

Optionally store channel delivery state in user_notifications.channel_state.

Suggested bits:

- bit 0: queued for SSE
- bit 1: delivered by SSE
- bit 2: queued for email
- bit 3: delivered by email
- bit 4: email failed
- bit 5: reserved for push queued
- bit 6: reserved for push delivered

This is separate from users.user_config.

users.user_config answers:
- what the user wants

user_notifications.channel_state answers:
- what happened for one notification instance

## 11. Package Layout

Suggested new packages based on current conventions:

```text
internal/pkg/model/notification.go
internal/pkg/model/user_notification.go
internal/pkg/model/user_config.go

internal/pkg/repo/notification/repo.go
internal/pkg/repo/notification/fx.go
internal/pkg/repo/user_notification/repo.go
internal/pkg/repo/user_notification/fx.go

internal/pkg/services/notification/service.go
internal/pkg/services/notification/fx.go
internal/pkg/services/notification/create_event.go
internal/pkg/services/notification/fanout.go
internal/pkg/services/notification/list.go
internal/pkg/services/notification/mark_seen.go
internal/pkg/services/notification/mark_read.go
internal/pkg/services/notification/mark_all_read.go
internal/pkg/services/notification/stream_publish.go

internal/app/api/route/notification/handler.go
internal/app/api/route/notification/route.go
internal/app/api/route/notification/fx.go
internal/app/api/route/notification/get_notifications.go
internal/app/api/route/notification/read_notification.go
internal/app/api/route/notification/read_all_notifications.go
internal/app/api/route/notification/stream_notifications.go

internal/pkg/notification/broker.go
internal/pkg/notification/payload.go
internal/pkg/notification/types.go
```

Suggested changes to existing modules:

- add notification repo modules to internal/pkg/repo/fx.go
- add notification service module to internal/pkg/services/fx.go
- add notification route module to internal/app/api/route/fx.go

## 12. User Config Type Design

Suggested implementation in Go:

```go
type UserConfigFlags uint64

func (f *UserConfigFlags) Scan(value any) error
func (f UserConfigFlags) Value() (driver.Value, error)
```

Storage note:

- GORM column type should be bigint
- default should be 0

Model update for user:

```go
type User struct {
	common.SqlModel
	Name string `json:"name" gorm:"column:name"`
	...
	UserConfig UserConfigFlags `json:"userConfig" gorm:"column:user_config;type:bytea;default:0"`
}
```

## 13. API Design

### 13.1 Notification APIs

Suggested route group:

- GET /notifications
- GET /notifications/stream
- PATCH /notifications/:id/seen
- PATCH /notifications/:id/read
- PATCH /notifications/read-all

Optional aggregate endpoint:

- GET /notifications/unread-count

Suggested list query params:

- page
- limit
- unreadOnly
- type

Suggested list item response:

```json
{
  "id": "uuid",
  "type": "comic.new_chapter",
  "title": "New chapter available",
  "body": "Chapter 12 of One Piece is now available",
  "isSeen": false,
  "isRead": false,
  "createdAt": "2026-04-15T10:00:00Z",
  "payload": {
    "comicId": "uuid",
    "comicSlug": "one-piece",
    "comicTitle": "One Piece",
    "chapterId": "uuid",
    "chapterSlug": "chapter-12",
    "chapterTitle": "Chapter 12"
  }
}
```

### 13.2 User Config APIs

Suggested endpoints:

- GET /users/me/config
- PATCH /users/me/config

PATCH request example:

```json
{
  "enableSseNotifications": true,
  "enableEmailNotifications": false,
  "enableComicNewChapterNotifications": true,
  "seenNotificationCenter": true
}
```

Behavior rules:

- PATCH only updates provided fields
- user_config bitmap is rebuilt by setting or clearing relevant bits
- response returns expanded JSON booleans, not raw bitmap

## 14. SSE Design

### 14.1 Why SSE

- one-way server-to-client fits notification delivery
- simpler than WebSocket for this use case
- native browser EventSource support
- easier infra footprint for initial rollout

### 14.2 Stream Endpoint

Suggested endpoint:

- GET /notifications/stream

Requirements:

- authenticated through existing JWT cookie middleware
- content type text/event-stream
- keep-alive heartbeat every 20 to 30 seconds
- graceful disconnect handling using request context

### 14.3 SSE Event Contract

Use named events.

Suggested event names:

- notification.created
- notification.badge_updated
- heartbeat

Example frame:

```text
id: <user_notification_id>
event: notification.created
data: {"id":"...","type":"comic.new_chapter",...}

```

### 14.4 Broker Design

The API layer should not fan out directly in memory only.

Recommended design:

- Redis Pub/Sub channel per user: notifications:user:{user_id}
- notification service publishes serialized payload to Redis
- SSE handler subscribes to the current user channel
- API instance forwards incoming pubsub messages to the HTTP stream

Why Redis Pub/Sub instead of only in-memory hub:

- works across multiple app instances
- no sticky session requirement
- reuses existing Redis dependency

### 14.5 SSE Delivery Rule

SSE should be treated as best effort real-time delivery.

Persistence rule:

- notification must be inserted into user_notifications before SSE publish

Client recovery rule:

- if SSE disconnects, client reconnects and fetches GET /notifications to fill gaps

## 15. Email Delivery Design

Email remains asynchronous and optional.

Rule evaluation order for comic.new_chapter:

1. user has enable email notifications bit
2. user has enable comic new chapter notifications bit
3. notification type supports email delivery

Suggested email process:

1. fan-out worker collects email-eligible recipients
2. worker batches recipients into a dedicated notification email task
3. notification email task builds type-specific mailables
4. existing mail system dispatches through mail queue

Email policy for phase 1:

- immediate delivery is acceptable
- dedupe by notification event id and user id
- skip email if recipient has no email or account is invalid

Future enhancement:

- digest emails every N minutes
- do not email if user was recently active through SSE

## 16. Asynq Task Design

Suggested new tasks:

- notification_fanout
- notification_email_delivery

Suggested payload for notification_fanout:

```json
{
  "type": "comic.new_chapter",
  "entityType": "chapter",
  "entityId": "<chapter_id>",
  "dedupeKey": "chapter-published:<chapter_id>",
  "triggeredBy": "<user_id>"
}
```

Suggested payload for notification_email_delivery:

```json
{
  "notificationId": "<notification_id>",
  "userNotificationIds": ["...", "..."],
  "recipientUserIds": ["...", "..."]
}
```

Guidelines:

- fan-out task is responsible for persistence and SSE publish
- email task is responsible only for email
- both tasks must be idempotent

## 17. Chapter Publish Integration

The first producer is chapter publish.

Trigger rule:

- only when chapter transitions from is_published = false to is_published = true

Do not notify when:

- chapter is created as draft
- chapter content is edited but remains published
- chapter is already published and publish API is called again

Recommended service change:

1. load chapter current state
2. if already published and request asks publish again, return success without enqueueing notification
3. if transitioning to published, perform DB update
4. enqueue notification_fanout task after successful update

Payload-specific data for comic.new_chapter should include:

- comic id
- comic slug
- comic title
- chapter id
- chapter slug
- chapter title
- chapter number

## 18. Recipient Resolution

For comic.new_chapter, recipients are active followers in comic_follows.

Query rules:

- only non-deleted comic_follows
- join users to access email and user_config
- skip the actor if product rule says creators should not receive their own notification

Future extensibility:

- each notification type can have its own recipient resolver
- resolvers can be registered by type

Suggested abstraction:

```go
type RecipientResolver interface {
	ResolveRecipients(ctx context.Context, event NotificationEventInput) ([]Recipient, error)
}
```

For phase 1, this can remain internal to NotificationService.

## 19. Idempotency and Consistency

This is critical.

### 19.1 Event Idempotency

Use notifications.dedupe_key.

For chapter published:

- dedupe_key = chapter-published:{chapter_id}

Behavior:

- if task retries, the worker should fetch-or-create by dedupe_key
- never create duplicate notification events for the same publication event

### 19.2 Recipient Idempotency

Use unique index on user_notifications(notification_id, user_id).

Behavior:

- insert with upsert or conflict do nothing
- retries must not duplicate inbox entries

### 19.3 Transaction Boundary

Recommended write order inside worker:

1. ensure notification event exists
2. insert recipient rows
3. commit transaction
4. publish SSE
5. enqueue email

Reasoning:

- user must be able to fetch notification even if SSE fails
- email should not be sent for a notification that was not persisted

## 20. Read and Seen Semantics

Define these clearly.

- seen: notification has been surfaced to the user interface
- read: user explicitly opened or acknowledged it

Suggested behavior:

- GET /notifications does not automatically mark items as seen
- PATCH /notifications/:id/seen marks a single item seen
- PATCH /notifications/:id/read marks a single item read and also seen
- PATCH /notifications/read-all marks all unread items as read and seen

This avoids hidden server-side side effects during list fetch.

## 21. Authorization Rules

- all notification endpoints require authenticated user
- list, read, seen, and stream endpoints always scope to current user only
- no endpoint should allow reading another user's notifications by ID alone

## 22. Observability

Recommended logs and metrics:

- notification events created by type
- recipient rows created by type
- SSE publishes by type
- SSE publish failures
- email queued by type
- email send failures
- unread count lookup latency

Tracing:

- include notification type, notification id, and user id in spans where applicable

## 23. Failure Handling

### 23.1 SSE Failure

- do not retry endlessly at producer side
- rely on persisted notification + client reconnect + list fetch

### 23.2 Email Failure

- retry through Asynq
- record failure in channel_state or log table if needed later

### 23.3 Partial Fan-Out Failure

- recipient insert should use batch operations
- if one batch fails, task retries safely because of idempotent constraints

## 24. Performance Considerations

- use batch insert for user_notifications
- paginate notification list by created_at desc
- do not preload excessive relations in inbox query
- keep payload concise and consumer-ready
- avoid per-user synchronous mail building in request path

For large follower counts:

- split recipients into chunks
- chunk email enqueueing
- consider async batch size limits for GORM insert and queue payload size

## 25. API Response Shape Strategy

Clients should not need to decode raw payload semantics from scratch.

Recommended list response structure:

- top-level generic fields: id, type, title, body, isSeen, isRead, createdAt
- payload carries navigation data and context-specific metadata

This allows a generic inbox UI while preserving type-specific actions.

## 26. Rollout Plan

### Phase 1

- add users.user_config
- add notifications and user_notifications tables
- add notification service and repos
- create list, seen, read, read-all APIs
- integrate chapter publish producer
- enqueue fan-out job

### Phase 2

- add Redis-backed SSE stream endpoint
- publish real-time events after persistence
- add unread count endpoint if needed

### Phase 3

- add email delivery policy for comic.new_chapter
- add notification-specific mail templates
- support more notification types

## 27. Open Product Decisions

These should be finalized before implementation:

1. Should creators receive notifications for their own chapter publish events?
2. Should unpublish then publish again create a second event?
3. Should opening the notification center auto-mark items as seen?
4. Should email delivery be immediate or digest-based in phase 1?
5. Should system announcements support all users or targeted groups?

## 28. Recommended Initial Defaults

For new users, set these defaults:

- enable SSE notifications: true
- enable email notifications: false
- enable comic new chapter notifications: true
- seen notification center: false

This keeps real-time delivery enabled by default while preventing surprise email volume.

## 29. Summary

The recommended design uses:

- one generic notifications table for source events
- one user_notifications table for recipient-specific state
- one bitmap-backed users.user_config field for stable user preferences
- Redis Pub/Sub for SSE fan-out across instances
- Asynq for background fan-out and email delivery

This design is generic enough to support future notification types without schema redesign, while keeping user settings compact and stored directly on the users table as requested.
