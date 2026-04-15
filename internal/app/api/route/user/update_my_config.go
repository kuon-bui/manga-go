package userroute

import (
	"manga-go/internal/app/api/common/response"
	userrequest "manga-go/internal/pkg/request/user"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Update current user config
// @Description  Update notification-related configuration for the current user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body  body      userrequest.UpdateUserConfigRequest  true  "User config update request"
// @Success      200   {object}  response.Result
// @Failure      400   {object}  response.Result
// @Failure      401   {object}  response.Result
// @Failure      500   {object}  response.Result
// @Router       /users/me/config [patch]
// @Security     AccessToken
func (h *userHandler) updateMyConfig(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	var req userrequest.UpdateUserConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.notificationService.UpdateUserConfig(c.Request.Context(), user.ID, &req)
	result.ResponseResult(c)
}
