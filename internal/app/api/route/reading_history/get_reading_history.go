package readinghistoryroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Get reading history by ID
// @Description  Retrieve a specific reading history record by its ID
// @Tags         ReadingHistory
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Reading history ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Security     AccessToken
// @Router       /reading-histories/{id} [get]
func (h *ReadingHistoryHandler) getReadingHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("Invalid id").ResponseResult(c)
		return
	}

	result := h.readingHistoryService.GetReadingHistory(c.Request.Context(), id)
	result.ResponseResult(c)
}
