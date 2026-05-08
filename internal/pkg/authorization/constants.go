package authorization

import (
	"fmt"

	"github.com/google/uuid"
)

type Org string
type Action string
type Object string
type Context string

const (
	OrgPlatform Org = "platform"

	// Contexts
	CtxAny         Context = "any"
	CtxOwner       Context = "owner"
	CtxGroupMember Context = "group_member"
	CtxGroupOwner  Context = "group_owner"
	CtxPublished   Context = "published"
	CtxSelf        Context = "self"

	// Actions
	ActionRead    Action = "read"
	ActionCreate  Action = "create"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionPublish Action = "publish"
	ActionManage  Action = "manage"

	// Objects
	ObjectAny              Object = "*"
	ObjectAuthor           Object = "author"
	ObjectChapter          Object = "chapter"
	ObjectComic            Object = "comic"
	ObjectComment          Object = "comment"
	ObjectFile             Object = "file"
	ObjectGenre            Object = "genre"
	ObjectNotification     Object = "notification"
	ObjectPage             Object = "page"
	ObjectPermission       Object = "permission"
	ObjectRating           Object = "rating"
	ObjectReadingHistory   Object = "reading_history"
	ObjectRole             Object = "role"
	ObjectTag              Object = "tag"
	ObjectTranslationGroup Object = "translation_group"
	ObjectUser             Object = "user"
)

func Subject(id uuid.UUID) string {
	return id.String()
}

func PlatformOrg() Org {
	return OrgPlatform
}

func TranslationGroupOrg(id uuid.UUID) Org {
	return Org(fmt.Sprintf("tg:%s", id.String()))
}

func TranslationGroupOrgString(id string) Org {
	return Org(fmt.Sprintf("tg:%s", id))
}

func DefaultContexts() []Context {
	return []Context{
		CtxOwner,
		CtxGroupOwner,
		CtxGroupMember,
		CtxPublished,
		CtxSelf,
		CtxAny,
	}
}
