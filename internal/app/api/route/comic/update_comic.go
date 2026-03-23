package comicroute

import (
	"manga-go/internal/app/api/common/response"
	comicrequest "manga-go/internal/pkg/request/comic"

	"github.com/gin-gonic/gin"
)

func (h *ComicHandler) updateComic(c *gin.Context) {
	slug := c.Param("slug")

	var req comicrequest.UpdateComicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.comicService.UpdateComic(c.Request.Context(), slug, &req)
	result.ResponseResult(c)
}
