package comicroute

import (
	"manga-go/internal/app/api/common/response"
	comicrequest "manga-go/internal/pkg/request/comic"

	"github.com/gin-gonic/gin"
)

// @Summary      Update comic
// @Description  Update comic information
// @Tags         Comic
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string                       true  "Comic slug"
// @Param        body       body      comicrequest.UpdateComicRequest  true  "Comic update request"
// @Success      200        {object}  response.Response
// @Failure      400        {object}  response.Response
// @Failure      401        {object}  response.Response
// @Failure      404        {object}  response.Response
// @Router       /comics/{comicSlug} [put]
// @Security     AccessToken
func (h *ComicHandler) updateComic(c *gin.Context) {
	slug := c.Param("comicSlug")

	var req comicrequest.UpdateComicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.comicService.UpdateComic(c.Request.Context(), slug, &req)
	result.ResponseResult(c)
}
