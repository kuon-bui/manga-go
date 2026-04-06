package authorroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      List authors
// @Description  Get paginated list of authors
// @Tags         Author
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Router       /authors [get]
// @Security     AccessToken
func (h *AuthorHandler) getAuthors(c *gin.Context) {
	var paging common.Paging
	if err := c.ShouldBindQuery(&paging); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.authorService.ListAuthors(c.Request.Context(), &paging)
	result.ResponseResult(c)
}
