package comicroute

import (
	_ "manga-go/internal/app/api/common/response"
	"github.com/gin-gonic/gin"
)

// @Summary      Get comic by slug
// @Description  Retrieve a specific comic by its slug
// @Tags         Comic
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string  true  "Comic slug"
// @Success      200        {object}  response.Response
// @Failure      400        {object}  response.Response
// @Failure      401        {object}  response.Response
// @Failure      404        {object}  response.Response
// @Router       /comics/{comicSlug} [get]
// @Security     AccessToken
func (h *ComicHandler) getComic(c *gin.Context) {
	slug := c.Param("comicSlug")

	result := h.comicService.GetComic(c.Request.Context(), slug)
	result.ResponseResult(c)
}
