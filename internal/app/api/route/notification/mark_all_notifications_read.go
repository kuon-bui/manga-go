package notificationroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Mark all notifications read
// @Description  Mark all notifications as read for the current user
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Result
// @Failure      401  {object}  response.Result
// @Failure      500  {object}  response.Result
// @Router       /notifications/read-all [patch]
// @Security     AccessToken
func (h *NotificationHandler) markAllNotificationsRead(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	result := h.notificationService.MarkAllNotificationsRead(c.Request.Context(), user.ID)
	result.ResponseResult(c)
}
