package permissionroute

import (
	"manga-go/internal/app/api/common/response"
	permissionrequest "manga-go/internal/pkg/request/permission"

	"github.com/gin-gonic/gin"
)

func (h *PermissionHandler) createPermission(c *gin.Context) {
	var req permissionrequest.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.permissionService.CreatePermission(c.Request.Context(), &req)
	result.ResponseResult(c)
}
