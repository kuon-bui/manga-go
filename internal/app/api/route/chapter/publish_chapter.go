package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"github.com/gin-gonic/gin"
)

// @Summary      Publish chapter
// @Description  Publish a chapter making it publicly visible
// @Tags         Chapter
// @Accept       json
// @Produce      json
// @Param        comicSlug    path      string                          true  "Comic slug"
// @Param        chapterSlug  path      string                          true  "Chapter slug"
// @Param        body         body      chapterrequest.PublishChapterRequest  true  "Publish request"
// @Success      200          {object}  response.Response
// @Failure      400          {object}  response.Response
// @Failure      401          {object}  response.Response
// @Failure      404          {object}  response.Response
// @Router       /comics/{comicSlug}/chapters/{chapterSlug}/publish [patch]
// @Security     AccessToken
func (h *ChapterHandler) publishChapter(c *gin.Context) {
	chapterSlug := c.Param("chapterSlug")

	var req chapterrequest.PublishChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.chapterService.PublishChapter(c.Request.Context(), chapterSlug, &req)
	result.ResponseResult(c)
}
