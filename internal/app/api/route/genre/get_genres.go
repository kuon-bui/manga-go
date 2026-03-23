package genreroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

func (h *GenreHandler) getGenres(c *gin.Context) {
	var paging common.Paging
	if err := c.ShouldBindQuery(&paging); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.genreService.ListGenres(c.Request.Context(), &paging)
	result.ResponseResult(c)
}
