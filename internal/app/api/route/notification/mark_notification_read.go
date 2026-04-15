package notificationroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Mark notification read
// @Description  Mark one notification as read for the current user
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Notification id"
// @Success      200  {object}  response.Result
// @Failure      400  {object}  response.Result
// @Failure      401  {object}  response.Result
// @Failure      500  {object}  response.Result
// @Router       /notifications/{id}/read [patch]
// @Security     AccessToken
func (h *NotificationHandler) markNotificationRead(c *gin.Context) {
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

	result := h.notificationService.MarkNotificationRead(c.Request.Context(), user.ID, id)
	result.ResponseResult(c)
}
