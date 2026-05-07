package authorization

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	OrgPlatform = "platform"

	CtxAny         = "any"
	CtxOwner       = "owner"
	CtxGroupMember = "group_member"
	CtxGroupOwner  = "group_owner"
	CtxPublished   = "published"
	CtxSelf        = "self"

	ActionRead    = "read"
	ActionCreate  = "create"
	ActionUpdate  = "update"
	ActionDelete  = "delete"
	ActionPublish = "publish"
	ActionManage  = "manage"

	ObjectAny              = "*"
	ObjectAuthor           = "author"
	ObjectChapter          = "chapter"
	ObjectComic            = "comic"
	ObjectComment          = "comment"
	ObjectFile             = "file"
	ObjectGenre            = "genre"
	ObjectNotification     = "notification"
	ObjectPage             = "page"
	ObjectPermission       = "permission"
	ObjectRating           = "rating"
	ObjectReadingHistory   = "reading_history"
	ObjectRole             = "role"
	ObjectTag              = "tag"
	ObjectTranslationGroup = "translation_group"
	ObjectUser             = "user"
)

func Subject(id uuid.UUID) string {
	return id.String()
}

func PlatformOrg() string {
	return OrgPlatform
}

func TranslationGroupOrg(id uuid.UUID) string {
	return fmt.Sprintf("tg:%s", id.String())
}

func DefaultContexts() []string {
	return []string{
		CtxOwner,
		CtxGroupOwner,
		CtxGroupMember,
		CtxPublished,
		CtxSelf,
		CtxAny,
	}
}
