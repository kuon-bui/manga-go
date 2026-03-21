package genreroute

import (
	"manga-go/internal/app/api/common/response"
	genrerequest "manga-go/internal/pkg/request/genre"

	"github.com/gin-gonic/gin"
)

func (h *GenreHandler) createGenre(c *gin.Context) {
	var req genrerequest.CreateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.genreService.CreateGenre(c.Request.Context(), &req)
	result.ResponseResult(c)
}
