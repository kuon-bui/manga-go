package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"github.com/gin-gonic/gin"
)

// @Summary      Update chapter pages
// @Description  Update pages content of a chapter
// @Tags         Chapter
// @Accept       json
// @Produce      json
// @Param        comicSlug    path      string                           true  "Comic slug"
// @Param        chapterSlug  path      string                           true  "Chapter slug"
// @Param        body         body      chapterrequest.UpdateChapterPagesRequest  true  "Pages update request"
// @Success      200          {object}  response.Response
// @Failure      400          {object}  response.Response
// @Failure      401          {object}  response.Response
// @Failure      404          {object}  response.Response
// @Router       /comics/{comicSlug}/chapters/{chapterSlug}/pages [put]
// @Security     AccessToken
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
