package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"github.com/gin-gonic/gin"
)

func (h *ChapterHandler) updateChapterPages(c *gin.Context) {
	chapterSlug := c.Param("chapterSlug")

	var req chapterrequest.UpdateChapterPagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.chapterService.UpdateChapterPages(c.Request.Context(), chapterSlug, &req)
	result.ResponseResult(c)
}
