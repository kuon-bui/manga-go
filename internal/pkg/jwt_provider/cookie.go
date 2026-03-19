package jwtprovider

import (
	"manga-go/internal/pkg/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SetCookieParams struct {
	AccessToken              string
	RefreshAccessToken       string
	ExpireAccessToken        time.Time
	ExpireRefreshAccessToken time.Time
}

func SetCookie(config *config.Config, g *gin.Context, p SetCookieParams) {
	isProduction := config.Production
	secure := isProduction
	sameSite := http.SameSiteLaxMode
	domain := ""

	if isProduction {
		sameSite = http.SameSiteNoneMode
		domain = config.Service.Domain
	}

	g.SetCookieData(&http.Cookie{
		Name:     config.CookieName.AccessToken,
		Value:    p.AccessToken,
		MaxAge:   int(time.Until(p.ExpireAccessToken).Seconds()),
		Path:     "/",
		Domain:   domain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	})
	g.SetCookieData(&http.Cookie{
		Name:     config.CookieName.RefreshToken,
		Value:    p.RefreshAccessToken,
		MaxAge:   int(time.Until(p.ExpireRefreshAccessToken).Seconds()),
		Path:     "/",
		Domain:   domain,
		Secure:   secure,
		HttpOnly: true,
		SameSite: sameSite,
	})
}
