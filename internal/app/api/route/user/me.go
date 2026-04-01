package userroute

import (
	_ "manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Get current user profile
// @Description  Retrieve the profile of the currently authenticated user
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /users/me [get]
// @Security     AccessToken
func (h *userHandler) me(c *gin.Context) {
	user, _ := utils.GetCurrentUserFromGinContext(c)
	c.JSON(200, gin.H{
		"user": user,
	})
}
