package genreroute

import (
	"manga-go/internal/app/api/common/response"
	genrerequest "manga-go/internal/pkg/request/genre"

	"github.com/gin-gonic/gin"
)

// @Summary      Update genre
// @Description  Update genre information
// @Tags         Genre
// @Accept       json
// @Produce      json
// @Param        slug  path      string                       true  "Genre slug"
// @Param        body  body      genrerequest.UpdateGenreRequest  true  "Genre update request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /genres/{slug} [put]
// @Security     AccessToken
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
