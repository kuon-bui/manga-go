package roleroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *RoleHandler) removeRolePermission(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid role id").ResponseResult(c)
		return
	}

	permissionID, err := uuid.Parse(c.Param("permissionId"))
	if err != nil {
		response.ResultError("invalid permission id").ResponseResult(c)
		return
	}

	result := h.roleService.RemovePermission(c.Request.Context(), roleID, permissionID)
	result.ResponseResult(c)
}
