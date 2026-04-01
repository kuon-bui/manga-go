package tagroute

import (
	_ "manga-go/internal/app/api/common/response"
	"github.com/gin-gonic/gin"
)

// @Summary      Get all tags
// @Description  Get all tags without pagination
// @Tags         Tag
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /tags/all [get]
// @Security     AccessToken
func (h *TagHandler) getAllTags(c *gin.Context) {
	result := h.tagService.ListAllTags(c.Request.Context())
	result.ResponseResult(c)
}
