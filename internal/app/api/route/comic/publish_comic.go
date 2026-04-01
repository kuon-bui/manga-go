package comicroute

import (
	"manga-go/internal/app/api/common/response"
	comicrequest "manga-go/internal/pkg/request/comic"

	"github.com/gin-gonic/gin"
)

// @Summary      Publish comic
// @Description  Publish a comic making it publicly visible
// @Tags         Comic
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string                        true  "Comic slug"
// @Param        body       body      comicrequest.PublishComicRequest  true  "Publish request"
// @Success      200        {object}  response.Response
// @Failure      400        {object}  response.Response
// @Failure      401        {object}  response.Response
// @Failure      404        {object}  response.Response
// @Router       /comics/{comicSlug}/publish [patch]
// @Security     AccessToken
func (h *ComicHandler) publishComic(c *gin.Context) {
	comicSlug := c.Param("comicSlug")

	var req comicrequest.PublishComicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.comicService.PublishComic(c.Request.Context(), comicSlug, &req)
	result.ResponseResult(c)
}
