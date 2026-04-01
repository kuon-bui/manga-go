package genreroute

import (
	_ "manga-go/internal/app/api/common/response"
	"github.com/gin-gonic/gin"
)

// @Summary      Get genre by slug
// @Description  Retrieve a specific genre by its slug
// @Tags         Genre
// @Accept       json
// @Produce      json
// @Param        slug  path      string  true  "Genre slug"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /genres/{slug} [get]
// @Security     AccessToken
func (h *GenreHandler) getGenre(c *gin.Context) {
	slug := c.Param("slug")

	result := h.genreService.GetGenre(c.Request.Context(), slug)
	result.ResponseResult(c)
}
