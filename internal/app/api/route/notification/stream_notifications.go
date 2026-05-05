package notificationroute

import (
	apicommon "manga-go/internal/app/api/common"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary      Stream notifications
// @Description  Stream real-time notifications for the current user using Server-Sent Events
// @Tags         Notification
// @Produce      text/event-stream
// @Success      200  {string}  string  "SSE stream"
// @Failure      401  {object}  response.Result
// @Router       /notifications/stream [get]
// @Security     AccessToken
func (h *NotificationHandler) streamNotifications(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	pubsub := h.notificationService.SubscribeUserChannel(c.Request.Context(), user.ID)
	defer pubsub.Close()

	stream := apicommon.NewSSEStream(c)
	if err := stream.SendComment("connected"); err != nil {
		return
	}

	heartbeat := time.NewTicker(25 * time.Second)
	defer heartbeat.Stop()

	channel := pubsub.Channel()
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-heartbeat.C:
			if err := stream.SendHeartbeat(); err != nil {
				return
			}
		case msg, ok := <-channel:
			if !ok {
				return
			}

			if err := stream.SendEvent("notification.created", msg.Payload); err != nil {
				return
			}
		}
	}
}
