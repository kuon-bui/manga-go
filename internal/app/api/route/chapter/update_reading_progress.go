package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	chapterrequest "manga-go/internal/pkg/request/chapter"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Update reading progress
// @Description  Update reading progress for the current user in a chapter
// @Tags         Chapter
// @Accept       json
// @Produce      json
// @Param        comicSlug    path      string                                      true  "Comic slug"
// @Param        chapterSlug  path      string                                      true  "Chapter slug"
// @Param        body         body      chapterrequest.UpdateReadingProgressRequest  true  "Reading progress update request"
// @Success      200          {object}  response.Response
// @Failure      400          {object}  response.Response
// @Failure      401          {object}  response.Response
// @Failure      404          {object}  response.Response
// @Router       /comics/{comicSlug}/chapters/{chapterSlug}/reading-progress [patch]
// @Security     AccessToken
func (h *ChapterHandler) updateReadingProgress(c *gin.Context) {
	ctx := c.Request.Context()
	chapterId, _ := common.GetChapterIdFromContext(ctx)
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}

	var req chapterrequest.UpdateReadingProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.chapterService.UpdateReadingProgress(ctx, user, chapterId, &req)
	result.ResponseResult(c)
}
