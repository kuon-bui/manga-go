package genreroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      List genres
// @Description  Get paginated list of genres
// @Tags         Genre
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Success      200    {object}  response.PaginationResponse
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Router       /genres [get]
// @Security     AccessToken
func (h *GenreHandler) getGenres(c *gin.Context) {
	var paging common.Paging
	if err := c.ShouldBindQuery(&paging); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.genreService.ListGenres(c.Request.Context(), &paging)
	result.ResponseResult(c)
}
