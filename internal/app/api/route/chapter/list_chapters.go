package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      List chapters
// @Description  Get paginated list of chapters for a comic
// @Tags         Chapter
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string  true  "Comic slug"
// @Param        page       query     int     false "Page number"
// @Param        limit      query     int     false "Items per page"
// @Success      200        {object}  response.PaginationResponse
// @Failure      400        {object}  response.Response
// @Failure      401        {object}  response.Response
// @Router       /comics/{comicSlug}/chapters [get]
// @Security     AccessToken
func (h *ChapterHandler) listChapters(c *gin.Context) {
	var paging common.Paging
	if err := c.ShouldBindQuery(&paging); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.chapterService.ListChapters(c.Request.Context(), &paging)
	result.ResponseResult(c)
}
