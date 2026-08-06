package roleroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Remove permission from role
// @Description  Revoke a single permission from a role, named as in GET /permissions, e.g. "comic:write".
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id              path      string  true  "Role ID"
// @Param        permissionName  path      string  true  "Permission name, e.g. comic:write"
// @Success      200             {object}  response.Result
// @Failure      400             {object}  response.Result
// @Failure      401             {object}  response.Result
// @Failure      403             {object}  response.Result
// @Router       /roles/{id}/permissions/{permissionName} [delete]
// @Security     AccessToken
func (h *RoleHandler) removeRolePermission(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid role id").ResponseResult(c)
		return
	}

	result := h.roleService.RemovePermission(c.Request.Context(), roleID, c.Param("permissionName"))
	result.ResponseResult(c)
}
