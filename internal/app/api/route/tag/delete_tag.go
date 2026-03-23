package tagroute

import (
	"github.com/gin-gonic/gin"
)

func (h *TagHandler) deleteTag(c *gin.Context) {
	slug := c.Param("slug")

	result := h.tagService.DeleteTag(c.Request.Context(), slug)
	result.ResponseResult(c)
}
