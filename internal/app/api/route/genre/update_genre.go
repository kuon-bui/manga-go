package genreroute

import (
	"manga-go/internal/app/api/common/response"
	genrerequest "manga-go/internal/pkg/request/genre"

	"github.com/gin-gonic/gin"
)

func (h *GenreHandler) updateGenre(c *gin.Context) {
	slug := c.Param("slug")

	var req genrerequest.UpdateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.genreService.UpdateGenre(c.Request.Context(), slug, &req)
	result.ResponseResult(c)
}
