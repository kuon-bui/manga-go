package notificationroute

import (
	"fmt"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/utils"
	"net/http"
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

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		response.ResultErrInternal(fmt.Errorf("streaming unsupported")).ResponseResult(c)
		return
	}

	_, _ = c.Writer.Write([]byte(": connected\n\n"))
	flusher.Flush()

	heartbeat := time.NewTicker(25 * time.Second)
	defer heartbeat.Stop()

	channel := pubsub.Channel()
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-heartbeat.C:
			_, _ = c.Writer.Write([]byte("event: heartbeat\ndata: {}\n\n"))
			flusher.Flush()
		case msg, ok := <-channel:
			if !ok {
				return
			}

			_, _ = c.Writer.Write([]byte("event: notification.created\n"))
			_, _ = c.Writer.Write([]byte("data: " + msg.Payload + "\n\n"))
			flusher.Flush()
		}
	}
}
