package genreroute

import (
	"github.com/gin-gonic/gin"
)

func (h *GenreHandler) getAllGenres(c *gin.Context) {
	result := h.genreService.ListAllGenres(c.Request.Context())
	result.ResponseResult(c)
}
