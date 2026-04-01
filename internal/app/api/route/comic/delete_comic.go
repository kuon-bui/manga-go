package comicroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

// @Summary      Delete comic
// @Description  Delete a comic by slug
// @Tags         Comic
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string  true  "Comic slug"
// @Success      200        {object}  response.Result
// @Failure      400        {object}  response.Result
// @Failure      401        {object}  response.Result
// @Failure      404        {object}  response.Result
// @Router       /comics/{comicSlug} [delete]
// @Security     AccessToken
func (h *ComicHandler) deleteComic(c *gin.Context) {
	slug := c.Param("comicSlug")
	var result response.Result
	result = h.comicService.DeleteComic(c.Request.Context(), slug)
	result.ResponseResult(c)
}
