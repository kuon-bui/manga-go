package userroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Get current user config
// @Description  Retrieve notification-related configuration for the current user
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Result
// @Failure      401  {object}  response.Result
// @Failure      500  {object}  response.Result
// @Router       /users/me/config [get]
// @Security     AccessToken
func (h *userHandler) getMyConfig(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	result := h.notificationService.GetUserConfig(c.Request.Context(), user.ID)
	result.ResponseResult(c)
}
