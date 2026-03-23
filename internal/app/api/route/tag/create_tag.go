package tagroute

import (
	"manga-go/internal/app/api/common/response"
	tagrequest "manga-go/internal/pkg/request/tag"

	"github.com/gin-gonic/gin"
)

func (h *TagHandler) createTag(c *gin.Context) {
	var req tagrequest.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.tagService.CreateTag(c.Request.Context(), &req)
	result.ResponseResult(c)
}
