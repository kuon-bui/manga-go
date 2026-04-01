package roleroute

import (
	"manga-go/internal/app/api/common/response"
	rolerequest "manga-go/internal/pkg/request/role"

	"github.com/gin-gonic/gin"
)

// @Summary      Create role
// @Description  Create a new role
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        body  body      rolerequest.CreateRoleRequest  true  "Role creation request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /roles [post]
// @Security     AccessToken
func (h *RoleHandler) createRole(c *gin.Context) {
	var req rolerequest.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.roleService.CreateRole(c.Request.Context(), &req)
	result.ResponseResult(c)
}
