package genreroute

import (
	_ "manga-go/internal/app/api/common/response"
	"github.com/gin-gonic/gin"
)

// @Summary      Delete genre
// @Description  Delete a genre by slug
// @Tags         Genre
// @Accept       json
// @Produce      json
// @Param        slug  path      string  true  "Genre slug"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /genres/{slug} [delete]
// @Security     AccessToken
func (h *GenreHandler) deleteGenre(c *gin.Context) {
	slug := c.Param("slug")

	result := h.genreService.DeleteGenre(c.Request.Context(), slug)
	result.ResponseResult(c)
}
