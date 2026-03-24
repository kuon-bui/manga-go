package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"github.com/gin-gonic/gin"
)

func (h *ChapterHandler) createChapter(c *gin.Context) {
	var req chapterrequest.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.chapterService.CreateChapter(c.Request.Context(), &req)
	result.ResponseResult(c)
}
