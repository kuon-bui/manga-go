package authzmiddleware

import (
	"errors"

	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	commentrepo "manga-go/internal/pkg/repo/comment"
	ratingrepo "manga-go/internal/pkg/repo/rating"
	readinghistoryrepo "manga-go/internal/pkg/repo/reading_history"
	translationgrouprepo "manga-go/internal/pkg/repo/translation_group"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthzMiddleware struct {
	authorizer           *authorization.Authorizer
	logger               *logger.Logger
	comicRepo            *comicrepo.ComicRepo
	chapterRepo          *chapterrepo.ChapterRepo
	commentRepo          *commentrepo.CommentRepo
	ratingRepo           *ratingrepo.RatingRepo
	readingHistoryRepo   *readinghistoryrepo.ReadingHistoryRepo
	translationGroupRepo *translationgrouprepo.TranslationGroupRepo
}

type AuthzMiddlewareParams struct {
	fx.In

	Authorizer           *authorization.Authorizer
	Logger               *logger.Logger
	ComicRepo            *comicrepo.ComicRepo
	ChapterRepo          *chapterrepo.ChapterRepo
	CommentRepo          *commentrepo.CommentRepo
	RatingRepo           *ratingrepo.RatingRepo
	ReadingHistoryRepo   *readinghistoryrepo.ReadingHistoryRepo
	TranslationGroupRepo *translationgrouprepo.TranslationGroupRepo
}

type ResourceContext struct {
	Org      authorization.Org
	Contexts []authorization.Context
}

type ResourceResolver func(*gin.Context, *model.User) (ResourceContext, error)

type authzContext struct {
	org      authorization.Org
	contexts []authorization.Context
}

func NewAuthzMiddleware(p AuthzMiddlewareParams) *AuthzMiddleware {
	return &AuthzMiddleware{
		authorizer:           p.Authorizer,
		logger:               p.Logger,
		comicRepo:            p.ComicRepo,
		chapterRepo:          p.ChapterRepo,
		commentRepo:          p.CommentRepo,
		ratingRepo:           p.RatingRepo,
		readingHistoryRepo:   p.ReadingHistoryRepo,
		translationGroupRepo: p.TranslationGroupRepo,
	}
}

func Require(m *AuthzMiddleware, action authorization.Action, resource authorization.Object, resolvers ...ResourceResolver) gin.HandlerFunc {
	if m == nil {
		return func(c *gin.Context) { c.Next() }
	}
	return m.Require(action, resource, resolvers...)
}

func (m *AuthzMiddleware) Require(action authorization.Action, resource authorization.Object, resolvers ...ResourceResolver) gin.HandlerFunc {
	return func(c *gin.Context) {
		if m == nil || m.authorizer == nil {
			c.Next()
			return
		}

		user, err := utils.GetCurrentUserFromGinContext(c)
		if err != nil {
			response.ResponseUnauthorized(c)
			c.Abort()
			return
		}

		baseReq := authorization.Request{
			Subject: authorization.Subject(user.ID),
			Org:     authorization.PlatformOrg(),
			Action:  action,
			Object:  resource,
			Context: authorization.CtxAny,
		}

		err = m.authorizer.Enforce(c.Request.Context(), baseReq)
		if err == nil {
			c.Next()
			return
		}
		if !errors.Is(err, authorization.ErrForbidden) {
			response.ResultErrInternal(err).ResponseResult(c)
			c.Abort()
			return
		}
		if len(resolvers) == 0 {
			response.ResponseForbidden(c)
			c.Abort()
			return
		}

		for _, resolver := range resolvers {
			resolved, err := resolver(c, user)
			if err != nil {
				m.respondResolverError(c, err)
				return
			}
			for _, item := range normalizeResourceContext(resolved) {
				req := baseReq
				req.Org = item.org
				if err := m.authorizer.EnforceAny(c.Request.Context(), req, item.contexts); err == nil {
					c.Next()
					return
				} else if !errors.Is(err, authorization.ErrForbidden) {
					response.ResultErrInternal(err).ResponseResult(c)
					c.Abort()
					return
				}
			}
		}

		response.ResponseForbidden(c)
		c.Abort()
	}
}

func (m *AuthzMiddleware) Comic() ResourceResolver {
	return func(c *gin.Context, user *model.User) (ResourceContext, error) {
		if m == nil || m.comicRepo == nil {
			return ResourceContext{}, errResolverUnavailable
		}

		id, ok := common.GetComicIdFromContext(c.Request.Context())
		if !ok {
			return ResourceContext{}, errResourceIDMissing
		}

		comic, err := m.comicRepo.FindOne(c.Request.Context(), []any{
			clause.Eq{Column: "id", Value: id},
		}, nil)
		if err != nil {
			return ResourceContext{}, err
		}

		contexts := make([]authorization.Context, 0, 4)
		org := authorization.PlatformOrg()
		if comic.UploadedByID != nil && *comic.UploadedByID == user.ID {
			contexts = append(contexts, authorization.CtxOwner)
		}
		if comic.TranslationGroupID != nil && user.TranslationGroupID != nil && *comic.TranslationGroupID == *user.TranslationGroupID {
			org = authorization.TranslationGroupOrg(*comic.TranslationGroupID)
			contexts = append(contexts, authorization.CtxGroupMember)
		}
		if comic.IsPublished {
			contexts = append(contexts, authorization.CtxPublished)
		}
		return ResourceContext{Org: org, Contexts: contexts}, nil
	}
}

func (m *AuthzMiddleware) Chapter() ResourceResolver {
	return func(c *gin.Context, user *model.User) (ResourceContext, error) {
		if m == nil || m.chapterRepo == nil {
			return ResourceContext{}, errResolverUnavailable
		}

		id, ok := common.GetChapterIdFromContext(c.Request.Context())
		if !ok {
			return ResourceContext{}, errResourceIDMissing
		}

		chapter, err := m.chapterRepo.FindOne(c.Request.Context(), []any{
			clause.Eq{Column: "id", Value: id},
		}, nil)
		if err != nil {
			return ResourceContext{}, err
		}

		contexts := make([]authorization.Context, 0, 4)
		org := authorization.PlatformOrg()
		if chapter.UploadedByID != nil && *chapter.UploadedByID == user.ID {
			contexts = append(contexts, authorization.CtxOwner)
		}
		if chapter.IsPublished {
			contexts = append(contexts, authorization.CtxPublished)
		}

		return ResourceContext{Org: org, Contexts: contexts}, nil
	}
}

func (m *AuthzMiddleware) ComicGroupFromContext() ResourceResolver {
	return func(c *gin.Context, user *model.User) (ResourceContext, error) {
		if m == nil || m.comicRepo == nil {
			return ResourceContext{}, errResolverUnavailable
		}

		id, ok := common.GetComicIdFromContext(c.Request.Context())
		if !ok {
			return ResourceContext{}, errResourceIDMissing
		}

		comic, err := m.comicRepo.FindOne(c.Request.Context(), []any{
			clause.Eq{Column: "id", Value: id},
		}, nil)
		if err != nil {
			return ResourceContext{}, err
		}

		contexts := make([]authorization.Context, 0, 2)
		org := authorization.PlatformOrg()
		if comic.TranslationGroupID != nil && user.TranslationGroupID != nil && *comic.TranslationGroupID == *user.TranslationGroupID {
			org = authorization.TranslationGroupOrg(*comic.TranslationGroupID)
			contexts = append(contexts, authorization.CtxGroupMember)
		}
		return ResourceContext{Org: org, Contexts: contexts}, nil
	}
}

func (m *AuthzMiddleware) CurrentUserTranslationGroup() ResourceResolver {
	return func(c *gin.Context, user *model.User) (ResourceContext, error) {
		if user.TranslationGroupID == nil {
			return ResourceContext{}, nil
		}

		return ResourceContext{
			Org:      authorization.TranslationGroupOrg(*user.TranslationGroupID),
			Contexts: []authorization.Context{authorization.CtxGroupMember},
		}, nil
	}
}

func (m *AuthzMiddleware) TranslationGroup() ResourceResolver {
	return func(c *gin.Context, user *model.User) (ResourceContext, error) {
		if m == nil || m.translationGroupRepo == nil {
			return ResourceContext{}, errResolverUnavailable
		}

		id, ok := common.GetTranslationGroupIdFromContext(c.Request.Context())
		if !ok {
			return ResourceContext{}, errResourceIDMissing
		}

		group, err := m.translationGroupRepo.FindOne(c.Request.Context(), []any{
			clause.Eq{Column: "id", Value: id},
		}, nil)
		if err != nil {
			return ResourceContext{}, err
		}

		contexts := make([]authorization.Context, 0, 2)
		if group.OwnerID == user.ID {
			contexts = append(contexts, authorization.CtxGroupOwner)
		}
		if user.TranslationGroupID != nil && *user.TranslationGroupID == group.ID {
			contexts = append(contexts, authorization.CtxGroupMember)
		}
		return ResourceContext{Org: authorization.TranslationGroupOrg(group.ID), Contexts: contexts}, nil
	}
}

func (m *AuthzMiddleware) CommentParam(param string) ResourceResolver {
	return func(c *gin.Context, user *model.User) (ResourceContext, error) {
		if m == nil || m.commentRepo == nil {
			return ResourceContext{}, errResolverUnavailable
		}

		id, err := uuid.Parse(c.Param(param))
		if err != nil {
			return ResourceContext{}, errInvalidResourceID
		}

		comment, err := m.commentRepo.FindOne(c.Request.Context(), []any{
			clause.Eq{Column: "id", Value: id},
		}, nil)
		if err != nil {
			return ResourceContext{}, err
		}

		contexts := []authorization.Context{}
		if comment.UserId == user.ID {
			contexts = append(contexts, authorization.CtxOwner)
		}
		return ResourceContext{Org: authorization.PlatformOrg(), Contexts: contexts}, nil
	}
}

func (m *AuthzMiddleware) RatingParam(param string) ResourceResolver {
	return func(c *gin.Context, user *model.User) (ResourceContext, error) {
		if m == nil || m.ratingRepo == nil {
			return ResourceContext{}, errResolverUnavailable
		}

		id, err := uuid.Parse(c.Param(param))
		if err != nil {
			return ResourceContext{}, errInvalidResourceID
		}

		rating, err := m.ratingRepo.FindOne(c.Request.Context(), []any{
			clause.Eq{Column: "id", Value: id},
		}, nil)
		if err != nil {
			return ResourceContext{}, err
		}

		contexts := []authorization.Context{}
		if rating.UserId == user.ID {
			contexts = append(contexts, authorization.CtxOwner)
		}
		return ResourceContext{Org: authorization.PlatformOrg(), Contexts: contexts}, nil
	}
}

func (m *AuthzMiddleware) ReadingHistoryParam(param string) ResourceResolver {
	return func(c *gin.Context, user *model.User) (ResourceContext, error) {
		if m == nil || m.readingHistoryRepo == nil {
			return ResourceContext{}, errResolverUnavailable
		}

		id, err := uuid.Parse(c.Param(param))
		if err != nil {
			return ResourceContext{}, errInvalidResourceID
		}

		readingHistory, err := m.readingHistoryRepo.FindOne(c.Request.Context(), []any{
			clause.Eq{Column: "id", Value: id},
		}, nil)
		if err != nil {
			return ResourceContext{}, err
		}

		contexts := []authorization.Context{}
		if readingHistory.UserID == user.ID {
			contexts = append(contexts, authorization.CtxOwner)
		}
		return ResourceContext{Org: authorization.PlatformOrg(), Contexts: contexts}, nil
	}
}

func UserParam(param string) ResourceResolver {
	return func(c *gin.Context, user *model.User) (ResourceContext, error) {
		id, err := uuid.Parse(c.Param(param))
		if err != nil {
			return ResourceContext{}, errInvalidResourceID
		}

		contexts := []authorization.Context{}
		if id == user.ID {
			contexts = append(contexts, authorization.CtxSelf)
		}
		return ResourceContext{Org: authorization.PlatformOrg(), Contexts: contexts}, nil
	}
}

var (
	errInvalidResourceID   = errors.New("invalid resource id")
	errResourceIDMissing   = errors.New("resource id missing")
	errResolverUnavailable = errors.New("authorization resource resolver unavailable")
)

func normalizeResourceContext(resolved ResourceContext) []authzContext {
	org := resolved.Org
	if org == "" {
		org = authorization.PlatformOrg()
	}
	contexts := resolved.Contexts
	if len(contexts) == 0 {
		contexts = []authorization.Context{authorization.CtxAny}
	}

	items := make([]authzContext, 0, len(contexts))
	for _, ctx := range contexts {
		if ctx == "" || ctx == authorization.CtxAny {
			continue
		}
		items = append(items, authzContext{
			org:      org,
			contexts: []authorization.Context{ctx},
		})
	}
	return items
}

func (m *AuthzMiddleware) respondResolverError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.ResponseNotFound(c, "Resource")
		c.Abort()
		return
	}
	if errors.Is(err, errInvalidResourceID) || errors.Is(err, errResourceIDMissing) {
		response.ResultError(err.Error()).ResponseResult(c)
		c.Abort()
		return
	}
	if m != nil && m.logger != nil {
		m.logger.Error("Failed to resolve authorization resource", "error", err)
	}
	response.ResultErrInternal(err).ResponseResult(c)
	c.Abort()
}
