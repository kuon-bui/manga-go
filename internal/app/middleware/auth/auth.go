package authmiddleware

import (
	"base-go/internal/app/api/common/response"
	"base-go/internal/pkg/config"
	jwtprovider "base-go/internal/pkg/jwt_provider"
	userrepo "base-go/internal/pkg/repo/user"
	"base-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm/clause"
)

type AuthMiddleware struct {
	jwt      *jwtprovider.JwtProvider
	config   *config.Config
	userRepo *userrepo.UserRepository
}

type AuthMiddlewareParams struct {
	fx.In

	Jwt      *jwtprovider.JwtProvider
	Config   *config.Config
	UserRepo *userrepo.UserRepository
}

func NewAuthMiddleware(params AuthMiddlewareParams) *AuthMiddleware {
	return &AuthMiddleware{
		jwt:      params.Jwt,
		config:   params.Config,
		userRepo: params.UserRepo,
	}
}

func (h *AuthMiddleware) RequireJwt(c *gin.Context) {
	accessToken := h.extractTokenFromCookies(c, h.config.CookieName.AccessToken)
	if c.IsAborted() {
		return
	}

	userContext, tokenId, err := h.jwt.ValidateToken(c.Request.Context(), accessToken)
	if err != nil || userContext == nil {
		response.ResponseUnauthorized(c)
		c.Abort()
		return
	}

	user, err := h.userRepo.FindOne(c.Request.Context(), []any{
		clause.Eq{
			Column: "id",
			Value:  userContext.UserID,
		},
	}, nil)
	if err != nil {
		response.ResponseUnauthorized(c)
		c.Abort()
		return
	}

	utils.SetCurrentUserToGinContext(c, user)
	utils.SetTokenIdToGinContext(c, tokenId)
}

func (h *AuthMiddleware) InvalidateJwt(c *gin.Context) {
	tokenId, err := utils.GetTokenIdFromGinContext(c)
	if err != nil {
		response.ResponseUnauthorized(c)
		c.Abort()
		return
	}

	err = h.jwt.InvalidateToken(c.Request.Context(), tokenId)
	if err != nil {
		response.ResponseUnauthorized(c)
		c.Abort()
		return
	}
}

func (h *AuthMiddleware) RenewToken(c *gin.Context) {
	refreshToken := h.extractTokenFromCookies(c, h.config.CookieName.RefreshToken)
	if c.IsAborted() {
		return
	}

	newAccessToken, newRefreshToken, err := h.jwt.RenewAccessToken(c.Request.Context(), refreshToken)
	if err != nil {
		response.ResponseUnauthorized(c)
		c.Abort()
		return
	}

	jwtprovider.SetCookie(h.config, c, jwtprovider.SetCookieParams{
		AccessToken:              newAccessToken.TokenString,
		RefreshAccessToken:       newRefreshToken.TokenString,
		ExpireAccessToken:        newAccessToken.ExpiresAt,
		ExpireRefreshAccessToken: newRefreshToken.ExpiresAt,
	})
}

func (m *AuthMiddleware) extractTokenFromCookies(c *gin.Context, name string) string {
	cookie := c.Request.CookiesNamed(name)
	if len(cookie) == 0 {
		response.ResponseUnauthorized(c)
		c.Abort()
		return ""
	}

	token := cookie[0]
	if token == nil || token.Value == "" {
		response.ResponseUnauthorized(c)
		c.Abort()
		return ""
	}

	return token.Value
}
