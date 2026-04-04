package chapterhandler

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Get reading progress
// @Description  Retrieve reading progress for the current user in a comic
// @Tags         Chapter
// @Accept       json
// @Produce      json
// @Param        comicSlug    path      string  true  "Comic slug"
// @Param        chapterSlug  path      string  true  "Chapter slug"
// @Success      200          {object}  response.Result
// @Failure      401          {object}  response.Result
// @Failure      404          {object}  response.Result
// @Failure      500          {object}  response.Result
// @Router       /comics/{comicSlug}/chapters/{chapterSlug}/reading-progress [get]
// @Security     AccessToken
func (h *ChapterHandler) getReadingProgress(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	result := h.chapterService.GetReadingProgress(c.Request.Context(), user)
	result.ResponseResult(c)
}
