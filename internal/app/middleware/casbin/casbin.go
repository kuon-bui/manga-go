package casbinmiddleware

import (
	"manga-go/internal/app/api/common/response"
	casbinpkg "manga-go/internal/pkg/casbin"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type CasbinMiddleware struct {
	enforcer *casbinpkg.Enforcer
	logger   *logger.Logger
}

type CasbinMiddlewareParams struct {
	fx.In

	Enforcer *casbinpkg.Enforcer
	Logger   *logger.Logger
}

func NewCasbinMiddleware(p CasbinMiddlewareParams) *CasbinMiddleware {
	return &CasbinMiddleware{
		enforcer: p.Enforcer,
		logger:   p.Logger,
	}
}

// RequireGlobalPermission checks if the current user has a global permission.
// For unauthenticated requests, uses "anonymous" as the subject.
func (m *CasbinMiddleware) RequireGlobalPermission(obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := m.getSubject(c)
		allowed, err := m.enforcer.Enforce(subject, casbinpkg.GlobalDomain, obj, act)
		if err != nil {
			m.logger.Errorf("Casbin enforce error: %v", err)
			response.ResultForbidden().ResponseResult(c)
			c.Abort()
			return
		}
		if !allowed {
			response.ResultForbidden().ResponseResult(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireGroupPermission checks if the current user has permission within the comic's translation group.
// The group ID is read from the request context (set by the slug middleware).
func (m *CasbinMiddleware) RequireGroupPermission(obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := utils.GetCurrentUserFromGinContext(c)
		if err != nil || user == nil {
			response.ResultUnauthorized().ResponseResult(c)
			c.Abort()
			return
		}

		groupID, ok := common.GetTranslationGroupIdFromContext(c.Request.Context())
		if !ok {
			response.ResultForbidden().ResponseResult(c)
			c.Abort()
			return
		}

		allowed, err := m.enforcer.Enforce(user.ID.String(), groupID.String(), obj, act)
		if err != nil {
			m.logger.Errorf("Casbin enforce error for user %s in group %s: %v", user.ID, groupID, err)
			response.ResultForbidden().ResponseResult(c)
			c.Abort()
			return
		}
		if !allowed {
			response.ResultForbidden().ResponseResult(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

// getSubject returns the user ID string if authenticated, otherwise "anonymous".
func (m *CasbinMiddleware) getSubject(c *gin.Context) string {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil || user == nil {
		return "anonymous"
	}
	return user.ID.String()
}
