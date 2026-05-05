package notificationroute

import (
	"errors"
	"io"
	"manga-go/internal/app/api/common/response"
	notificationrequest "manga-go/internal/pkg/request/notification"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Create test notification
// @Description  Create a test notification for the current user and publish it over SSE
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param        body  body      notificationrequest.CreateTestNotificationRequest  false  "Test notification request"
// @Success      200   {object}  response.Result
// @Failure      400   {object}  response.Result
// @Failure      401   {object}  response.Result
// @Failure      500   {object}  response.Result
// @Router       /notifications/test [post]
// @Security     AccessToken
func (h *NotificationHandler) createTestNotification(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	var req notificationrequest.CreateTestNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.notificationService.CreateTestNotification(c.Request.Context(), user.ID, &req)
	result.ResponseResult(c)
}
