package readinghistoryroute

import (
	"manga-go/internal/app/api/common/response"
	readinghistoryrequest "manga-go/internal/pkg/request/reading_history"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Create reading history
// @Description  Create a new reading history record for current user
// @Tags         ReadingHistory
// @Accept       json
// @Produce      json
// @Param        body  body      readinghistoryrequest.CreateReadingHistoryRequest  true  "Reading history creation request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Security     AccessToken
// @Router       /reading-histories [post]
func (h *ReadingHistoryHandler) createReadingHistory(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	var req readinghistoryrequest.CreateReadingHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.readingHistoryService.CreateReadingHistory(c.Request.Context(), user.ID, &req)
	result.ResponseResult(c)
}
