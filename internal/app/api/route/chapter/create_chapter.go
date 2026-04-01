package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	chapterrequest "manga-go/internal/pkg/request/chapter"

	"github.com/gin-gonic/gin"
)

// @Summary      Create chapter
// @Description  Create a new chapter for a comic
// @Tags         Chapter
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string                        true  "Comic slug"
// @Param        body       body      chapterrequest.CreateChapterRequest  true  "Chapter creation request"
// @Success      200        {object}  response.Response
// @Failure      400        {object}  response.Response
// @Failure      401        {object}  response.Response
// @Router       /comics/{comicSlug}/chapters [post]
// @Security     AccessToken
func (h *ChapterHandler) createChapter(c *gin.Context) {
	var req chapterrequest.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.chapterService.CreateChapter(c.Request.Context(), &req)
	result.ResponseResult(c)
}
