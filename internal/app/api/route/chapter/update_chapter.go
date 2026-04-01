package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"github.com/gin-gonic/gin"
)

// @Summary      Update chapter
// @Description  Update chapter information
// @Tags         Chapter
// @Accept       json
// @Produce      json
// @Param        comicSlug    path      string                        true  "Comic slug"
// @Param        chapterSlug  path      string                        true  "Chapter slug"
// @Param        body         body      chapterrequest.UpdateChapterRequest  true  "Chapter update request"
// @Success      200          {object}  response.Response
// @Failure      400          {object}  response.Response
// @Failure      401          {object}  response.Response
// @Failure      404          {object}  response.Response
// @Router       /comics/{comicSlug}/chapters/{chapterSlug} [put]
// @Security     AccessToken
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
