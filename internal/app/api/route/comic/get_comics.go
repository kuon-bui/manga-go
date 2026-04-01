package comicroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      List comics
// @Description  Get paginated list of comics
// @Tags         Comic
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Success      200    {object}  response.PaginationResponse
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Router       /comics [get]
// @Security     AccessToken
func (h *ComicHandler) getComics(c *gin.Context) {
	var paging common.Paging
	if err := c.ShouldBindQuery(&paging); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.comicService.ListComics(c.Request.Context(), &paging)
	result.ResponseResult(c)
}
