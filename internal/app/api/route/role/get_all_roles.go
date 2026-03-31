package roleroute

import "github.com/gin-gonic/gin"

func (h *RoleHandler) getAllRoles(c *gin.Context) {
	result := h.roleService.ListAllRoles(c.Request.Context())
	result.ResponseResult(c)
}
