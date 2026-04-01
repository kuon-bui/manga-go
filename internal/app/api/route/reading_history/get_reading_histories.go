package readinghistoryroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      List reading histories
// @Description  Retrieve paginated reading history records for current user
// @Tags         ReadingHistory
// @Accept       json
// @Produce      json
// @Param        page   query  int  false  "Page number (default: 1)"
// @Param        limit  query  int  false  "Records per page (default: 10)"
// @Success      200    {object}  response.PaginationResponse
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Failure      500    {object}  response.Response
// @Security     AccessToken
// @Router       /reading-histories [get]
func (h *ReadingHistoryHandler) getReadingHistories(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	var paging common.Paging
	if err := c.ShouldBindQuery(&paging); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.readingHistoryService.ListReadingHistories(c.Request.Context(), user.ID, &paging)
	result.ResponseResult(c)
}
