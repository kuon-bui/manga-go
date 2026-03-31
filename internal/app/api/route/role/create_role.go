package roleroute

import (
	"manga-go/internal/app/api/common/response"
	rolerequest "manga-go/internal/pkg/request/role"

	"github.com/gin-gonic/gin"
)

func (h *RoleHandler) createRole(c *gin.Context) {
	var req rolerequest.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.roleService.CreateRole(c.Request.Context(), &req)
	result.ResponseResult(c)
}
