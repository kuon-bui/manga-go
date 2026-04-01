package userroute

import (
	_ "manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

// @Summary      Renew access token
// @Description  Renew expired access token using refresh token cookie. Returns new access_token and refresh_token cookies.
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Security     RefreshToken
// @Router       /users/renew-token [post]
func (h *userHandler) renewToken(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "renew token",
	})
}
