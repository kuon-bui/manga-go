package roleroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Remove permission from role
// @Description  Remove a permission from a role
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id            path      string  true  "Role ID"
// @Param        permissionId  path      string  true  "Permission ID"
// @Success      200           {object}  response.Response
// @Failure      400           {object}  response.Response
// @Failure      401           {object}  response.Response
// @Router       /roles/{id}/permissions/{permissionId} [delete]
// @Security     AccessToken
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
