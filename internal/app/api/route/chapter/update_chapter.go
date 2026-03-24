package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"github.com/gin-gonic/gin"
)

func (h *ChapterHandler) updateChapter(c *gin.Context) {
	chapterSlug := c.Param("chapterSlug")

	var req chapterrequest.UpdateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.chapterService.UpdateChapter(c.Request.Context(), chapterSlug, &req)
	result.ResponseResult(c)
}
