package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"github.com/gin-gonic/gin"
)

// @Summary      Get recent chapter updates
// @Description  Get list of recently updated chapters across all comics
// @Tags         Chapter
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Result
// @Failure      401    {object}  response.Result
// @Failure      500    {object}  response.Result
// @Security     AccessToken
// @Router       /chapters/recent-updates [get]
func (h *ChapterHandler) getRecentUpdates(c *gin.Context) {
	var req chapterrequest.RecentUpdatesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.chapterService.GetRecentUpdates(c.Request.Context(), &req)
	result.ResponseResult(c)
}
