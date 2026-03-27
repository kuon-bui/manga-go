package comicroute

import (
	"github.com/gin-gonic/gin"
)

func (h *ComicHandler) getComic(c *gin.Context) {
	slug := c.Param("comicSlug")

	result := h.comicService.GetComic(c.Request.Context(), slug)
	result.ResponseResult(c)
}
