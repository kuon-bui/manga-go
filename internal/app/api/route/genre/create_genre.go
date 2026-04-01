package genreroute

import (
	"manga-go/internal/app/api/common/response"
	genrerequest "manga-go/internal/pkg/request/genre"

	"github.com/gin-gonic/gin"
)

// @Summary      Create genre
// @Description  Create a new genre
// @Tags         Genre
// @Accept       json
// @Produce      json
// @Param        body  body      genrerequest.CreateGenreRequest  true  "Genre creation request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /genres [post]
// @Security     AccessToken
func (h *GenreHandler) createGenre(c *gin.Context) {
	var req genrerequest.CreateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.genreService.CreateGenre(c.Request.Context(), &req)
	result.ResponseResult(c)
}
