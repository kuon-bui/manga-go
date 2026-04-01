package userroute

import (
	"manga-go/internal/app/api/common/response"
	jwtprovider "manga-go/internal/pkg/jwt_provider"
	userrequest "manga-go/internal/pkg/request/user"

	"github.com/gin-gonic/gin"
)

// @Summary      Sign in user
// @Description  Authenticate user with email and password. On success, sets access_token and refresh_token HTTP-only cookies.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body  body      userrequest.SignInRequest  true  "Sign in request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /users/sign-in [post]
func (h *userHandler) signIn(c *gin.Context) {
	var req userrequest.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	accessToken, refreshToken, res := h.userService.SignIn(c.Request.Context(), &req)
	if accessToken != nil && refreshToken != nil {
		jwtprovider.SetCookie(h.config, c, jwtprovider.SetCookieParams{
			AccessToken:              accessToken.TokenString,
			RefreshAccessToken:       refreshToken.TokenString,
			ExpireAccessToken:        accessToken.ExpiresAt,
			ExpireRefreshAccessToken: refreshToken.ExpiresAt,
		})
	}

	res.ResponseResult(c)
}
