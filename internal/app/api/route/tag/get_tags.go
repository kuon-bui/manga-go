package tagroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      List tags
// @Description  Get paginated list of tags
// @Tags         Tag
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Router       /tags [get]
// @Security     AccessToken
func (h *TagHandler) getTags(c *gin.Context) {
	var paging common.Paging
	if err := c.ShouldBindQuery(&paging); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.tagService.ListTags(c.Request.Context(), &paging)
	result.ResponseResult(c)
}
