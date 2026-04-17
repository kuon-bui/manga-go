package chapterhandler

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

// @Summary      Get chapter by slug
// @Description  Retrieve a specific chapter by its slug
// @Tags         Chapter
// @Accept       json
// @Produce      json
// @Param        comicSlug    path      string  true  "Comic slug"
// @Param        chapterSlug  path      string  true  "Chapter slug"
// @Success      200          {object}  response.Result
// @Failure      400          {object}  response.Result
// @Failure      401          {object}  response.Result
// @Failure      404          {object}  response.Result
// @Router       /comics/{comicSlug}/chapters/{chapterSlug} [get]
// @Security     AccessToken
func (h *ChapterHandler) getChapter(c *gin.Context) {
	chapterSlug := c.Param("chapterSlug")
	var result response.Result
	result = h.chapterService.GetChapter(c.Request.Context(), chapterSlug)
	h.normalizeChapterImageURLs(c, &result)
	result.ResponseResult(c)
}
