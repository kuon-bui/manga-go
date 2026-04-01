package genreroute

import (
	_ "manga-go/internal/app/api/common/response"
	"github.com/gin-gonic/gin"
)

// @Summary      Get all genres
// @Description  Get all genres without pagination
// @Tags         Genre
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /genres/all [get]
// @Security     AccessToken
func (h *GenreHandler) getAllGenres(c *gin.Context) {
	result := h.genreService.ListAllGenres(c.Request.Context())
	result.ResponseResult(c)
}
