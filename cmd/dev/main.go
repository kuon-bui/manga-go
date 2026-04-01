package main

import (
	"manga-go/internal/app"
	"manga-go/internal/app/api"
	"manga-go/internal/app/api/server"
	"manga-go/internal/pkg/tracer"

	"go.uber.org/fx"
)

// @title           Manga-Go API
// @version         1.0
// @description     REST API for manga reading application.
// @description     Authentication uses HTTP-only cookies. After signing in via POST /users/sign-in,
// @description     the server sets two cookies: **access_token** (short-lived JWT) and **refresh_token** (long-lived JWT).
// @description     All protected endpoints require a valid access_token cookie.
// @description     Use the "Authorize" button below and paste your access_token value to test protected endpoints.
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey AccessToken
// @in             cookie
// @name           access_token
// @description    JWT access token stored in HTTP-only cookie. Obtained from POST /users/sign-in.
// @securityDefinitions.apikey RefreshToken
// @in             cookie
// @name           refresh_token
// @description    JWT refresh token stored in HTTP-only cookie. Used by POST /users/renew-token.
func main() {
	fx.New(
		app.Module,
		api.Module,
		tracer.Module,
		fx.Invoke(server.RunServer),
	).Run()
}
