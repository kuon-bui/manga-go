package permissionroute

import (
	"manga-go/internal/app/api/common/response"
	permissionrequest "manga-go/internal/pkg/request/permission"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *PermissionHandler) updatePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	var req permissionrequest.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.permissionService.UpdatePermission(c.Request.Context(), id, &req)
	result.ResponseResult(c)
}
