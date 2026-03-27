package tagroute

import (
	"github.com/gin-gonic/gin"
)

func (h *TagHandler) getAllTags(c *gin.Context) {
	result := h.tagService.ListAllTags(c.Request.Context())
	result.ResponseResult(c)
}
