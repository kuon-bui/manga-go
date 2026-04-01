package tagroute

import (
	_ "manga-go/internal/app/api/common/response"
	"github.com/gin-gonic/gin"
)

// @Summary      Delete tag
// @Description  Delete a tag by slug
// @Tags         Tag
// @Accept       json
// @Produce      json
// @Param        slug  path      string  true  "Tag slug"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /tags/{slug} [delete]
// @Security     AccessToken
func (h *TagHandler) deleteTag(c *gin.Context) {
	slug := c.Param("slug")

	result := h.tagService.DeleteTag(c.Request.Context(), slug)
	result.ResponseResult(c)
}
