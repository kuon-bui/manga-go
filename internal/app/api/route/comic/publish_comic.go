package comicroute

import (
	"manga-go/internal/app/api/common/response"
	comicrequest "manga-go/internal/pkg/request/comic"

	"github.com/gin-gonic/gin"
)

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
