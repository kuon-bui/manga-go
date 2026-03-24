package comicroute

import (
	"github.com/gin-gonic/gin"
)

func (h *ComicHandler) deleteComic(c *gin.Context) {
	slug := c.Param("comicSlug")

	result := h.comicService.DeleteComic(c.Request.Context(), slug)
	result.ResponseResult(c)
}
