package comicroute

import (
	"manga-go/internal/app/api/common/response"
	comicrequest "manga-go/internal/pkg/request/comic"

	"github.com/gin-gonic/gin"
)

func (h *ComicHandler) createComic(c *gin.Context) {
	var req comicrequest.CreateComicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.comicService.CreateComic(c.Request.Context(), &req)
	result.ResponseResult(c)
}
