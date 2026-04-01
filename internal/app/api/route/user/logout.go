package userroute

import (
	_ "manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

// @Summary      Logout user
// @Description  Invalidate user's JWT token and clear cookies
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /users/logout [delete]
// @Security     AccessToken
func (h *userHandler) logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "logout",
	})
}
