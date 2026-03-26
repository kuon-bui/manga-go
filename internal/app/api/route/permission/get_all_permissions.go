package permissionroute

import "github.com/gin-gonic/gin"

func (h *PermissionHandler) getAllPermissions(c *gin.Context) {
	result := h.permissionService.ListAllPermissions(c.Request.Context())
	result.ResponseResult(c)
}
