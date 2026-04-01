package readinghistoryroute

import (
	"manga-go/internal/app/api/common/response"
	readinghistoryrequest "manga-go/internal/pkg/request/reading_history"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Update reading history
// @Description  Update reading history record (typically to update last_read_at timestamp)
// @Tags         ReadingHistory
// @Accept       json
// @Produce      json
// @Param        id    path      string                                      true  "Reading history ID"
// @Param        body  body      readinghistoryrequest.UpdateReadingHistoryRequest  true  "Reading history update request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Security     AccessToken
// @Router       /reading-histories/{id} [put]
func (h *ReadingHistoryHandler) updateReadingHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("Invalid id").ResponseResult(c)
		return
	}

	var req readinghistoryrequest.UpdateReadingHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.readingHistoryService.UpdateReadingHistory(c.Request.Context(), id, &req)
	result.ResponseResult(c)
}
