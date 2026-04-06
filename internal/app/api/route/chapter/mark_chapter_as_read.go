package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      Mark chapter as read
// @Description  Mark a chapter as read for the current user
// @Tags         Chapter
// @Accept       json
// @Produce      json
// @Param        comicSlug    path      string  true  "Comic slug"
// @Param        chapterSlug  path      string  true  "Chapter slug"
// @Success      200          {object}  response.Response
// @Failure      400          {object}  response.Response
// @Failure      401          {object}  response.Response
// @Failure      404          {object}  response.Response
// @Failure      500          {object}  response.Response
// @Router       /comics/{comicSlug}/chapters/{chapterSlug}/mark-as-read [patch]
// @Security     AccessToken
func (h *ChapterHandler) markChapterAsRead(c *gin.Context) {
	ctx := c.Request.Context()

	comicID, ok := common.GetComicIdFromContext(ctx)
	if !ok {
		response.ResultError("Comic not found in context").ResponseResult(c)
		return
	}

	chapterID, ok := common.GetChapterIdFromContext(ctx)
	if !ok {
		response.ResultError("Chapter not found in context").ResponseResult(c)
		return
	}

	result := h.chapterService.MarkChapterAsRead(ctx, comicID, chapterID)
	result.ResponseResult(c)
}
