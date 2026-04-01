package tagroute

import (
	"manga-go/internal/app/api/common/response"
	tagrequest "manga-go/internal/pkg/request/tag"

	"github.com/gin-gonic/gin"
)

// @Summary      Create tag
// @Description  Create a new tag
// @Tags         Tag
// @Accept       json
// @Produce      json
// @Param        body  body      tagrequest.CreateTagRequest  true  "Tag creation request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /tags [post]
// @Security     AccessToken
func (h *TagHandler) createTag(c *gin.Context) {
	var req tagrequest.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.tagService.CreateTag(c.Request.Context(), &req)
	result.ResponseResult(c)
}
