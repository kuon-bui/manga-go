package roleroute

import (
	"strings"

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
// @Param        If-Match        header    string  false "Current global authorization version (for example g12)"
// @Success      200             {object}  response.Result
// @Failure      400             {object}  response.Result
// @Failure      401             {object}  response.Result
// @Failure      403             {object}  response.Result
// @Failure      404             {object}  response.Result
// @Failure      409             {object}  response.Result  "AUTHORIZATION_STATE_CHANGED, SELF_MANAGE_REQUIRED, or LAST_ROLE_MANAGER"
// @Router       /roles/{id}/permissions/{permissionName} [delete]
// @Security     AccessToken
func (h *RoleHandler) removeRolePermission(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid role id").ResponseResult(c)
		return
	}

	result := h.roleService.RemovePermission(
		c.Request.Context(), roleID, c.Param("permissionName"), roleAuthorizationVersion(c),
	)
	result.ResponseResult(c)
}

func roleAuthorizationVersion(c *gin.Context) string {
	return strings.Trim(c.GetHeader("If-Match"), "\"")
}
