package genreroute

import (
	"github.com/gin-gonic/gin"
)

func (h *GenreHandler) deleteGenre(c *gin.Context) {
	slug := c.Param("slug")

	result := h.genreService.DeleteGenre(c.Request.Context(), slug)
	result.ResponseResult(c)
}
