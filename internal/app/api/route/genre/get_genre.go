package genreroute

import (
	"github.com/gin-gonic/gin"
)

func (h *GenreHandler) getGenre(c *gin.Context) {
	slug := c.Param("slug")

	result := h.genreService.GetGenre(c.Request.Context(), slug)
	result.ResponseResult(c)
}
