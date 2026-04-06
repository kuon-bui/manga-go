package roleroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      List roles
// @Description  Get paginated list of roles
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Router       /roles [get]
// @Security     AccessToken
func (h *RoleHandler) getRoles(c *gin.Context) {
	var paging common.Paging
	if err := c.ShouldBindQuery(&paging); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.roleService.ListRoles(c.Request.Context(), &paging)
	result.ResponseResult(c)
}
