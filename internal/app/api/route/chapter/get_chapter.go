package chapterhandler

import (
	"github.com/gin-gonic/gin"
)

func (h *ChapterHandler) getChapter(c *gin.Context) {
	chapterSlug := c.Param("chapterSlug")
	result := h.chapterService.GetChapter(c.Request.Context(), chapterSlug)
	result.ResponseResult(c)
}
