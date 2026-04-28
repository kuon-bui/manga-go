package jwtprovider

import (
	"manga-go/internal/pkg/config"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SetCookieParams struct {
	AccessToken              string
	RefreshAccessToken       string
	ExpireAccessToken        time.Time
	ExpireRefreshAccessToken time.Time
}

func SetCookie(cfg *config.Config, g *gin.Context, p SetCookieParams) {
	isProduction := cfg.RunMode == config.RunModeProduction
	domain := ""

	// Detect if the request is cross-origin by checking the Origin header.
	// If the frontend is served from a different host (e.g. a tunnel like nport.link),
	// we must use SameSite=None + Secure=true so the browser will attach the cookie.
	// If the request comes from localhost/127.0.0.1, use SameSite=Lax + Secure=false
	// so the cookie works correctly over plain HTTP during local development.
	secure := false
	sameSite := http.SameSiteLaxMode

	origin := g.GetHeader("Origin")
	if origin != "" &&
		!strings.HasPrefix(origin, "http://localhost:") &&
		!strings.HasPrefix(origin, "https://localhost:") &&
		!strings.HasPrefix(origin, "http://127.0.0.1:") &&
		!strings.HasPrefix(origin, "https://127.0.0.1:") {
		// Cross-origin request: browser requires SameSite=None + Secure
		secure = true
		sameSite = http.SameSiteNoneMode
	}

	if isProduction {
		domain = cfg.Service.Domain
	}

	g.SetCookieData(&http.Cookie{
		Name:     cfg.CookieName.AccessToken,
		Value:    p.AccessToken,
		MaxAge:   int(time.Until(p.ExpireAccessToken).Seconds()),
		Path:     "/",
		Domain:   domain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	})
	g.SetCookieData(&http.Cookie{
		Name:     cfg.CookieName.RefreshToken,
		Value:    p.RefreshAccessToken,
		MaxAge:   int(time.Until(p.ExpireRefreshAccessToken).Seconds()),
		Path:     "/",
		Domain:   domain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	})
}
