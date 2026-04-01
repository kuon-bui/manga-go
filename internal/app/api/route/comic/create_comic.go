package comicroute

import (
	"manga-go/internal/app/api/common/response"
	comicrequest "manga-go/internal/pkg/request/comic"

	"github.com/gin-gonic/gin"
)

// @Summary      Create comic
// @Description  Create a new comic
// @Tags         Comic
// @Accept       json
// @Produce      json
// @Param        body  body      comicrequest.CreateComicRequest  true  "Comic creation request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /comics [post]
// @Security     AccessToken
func (h *ComicHandler) createComic(c *gin.Context) {
	var req comicrequest.CreateComicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.comicService.CreateComic(c.Request.Context(), &req)
	result.ResponseResult(c)
}
