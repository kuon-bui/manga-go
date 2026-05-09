package notificationroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Delete notification
// @Description  Soft delete a notification for the current user
// @Tags         Notification
// @Produce      json
// @Param        id   path      string  true  "Notification id"
// @Success      200  {object}  response.Result
// @Failure      400  {object}  response.Result
// @Failure      401  {object}  response.Result
// @Failure      404  {object}  response.Result
// @Failure      500  {object}  response.Result
// @Router       /notifications/{id} [delete]
// @Security     AccessToken
func (h *NotificationHandler) deleteNotification(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("Invalid id").ResponseResult(c)
		return
	}

	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	result := h.notificationService.DeleteNotification(c.Request.Context(), user.ID, id)
	result.ResponseResult(c)
}

// @Summary      Delete all notifications
// @Description  Soft delete all notifications for the current user
// @Tags         Notification
// @Produce      json
// @Success      200  {object}  response.Result
// @Failure      401  {object}  response.Result
// @Failure      500  {object}  response.Result
// @Router       /notifications [delete]
// @Security     AccessToken
func (h *NotificationHandler) deleteAllNotifications(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	result := h.notificationService.DeleteAllNotifications(c.Request.Context(), user.ID)
	result.ResponseResult(c)
}
