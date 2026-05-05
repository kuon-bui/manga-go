package notificationroute

import (
	"manga-go/internal/app/api/common/response"
	notificationrequest "manga-go/internal/pkg/request/notification"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      List notifications
// @Description  Retrieve paginated notifications for the current user
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "Page"
// @Param        limit       query     int     false  "Limit"
// @Param        unreadOnly  query     bool    false  "Unread only"
// @Param        type        query     string  false  "Notification type"
// @Success      200         {object}  response.Result
// @Failure      400         {object}  response.Result
// @Failure      401         {object}  response.Result
// @Failure      500         {object}  response.Result
// @Router       /notifications [get]
// @Security     AccessToken
func (h *NotificationHandler) getNotifications(c *gin.Context) {
	var req notificationrequest.ListNotificationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	result := h.notificationService.ListNotifications(c.Request.Context(), user.ID, &req)
	result.ResponseResult(c)
}
