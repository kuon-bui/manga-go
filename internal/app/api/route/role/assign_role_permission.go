package roleroute

import (
	"manga-go/internal/app/api/common/response"
	rolerequest "manga-go/internal/pkg/request/role"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Assign permissions to role
// @Description  Replace the role's permissions with the given set. Names come from GET /permissions, e.g. "comic:write".
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id    path      string                               true  "Role ID"
// @Param        body  body      rolerequest.AssignPermissionRequest   true  "Permissions to assign"
// @Success      200   {object}  response.Result
// @Failure      400   {object}  response.Result
// @Failure      401   {object}  response.Result
// @Failure      403   {object}  response.Result
// @Router       /roles/{id}/permissions [post]
// @Security     AccessToken
func (h *RoleHandler) assignRolePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	var req rolerequest.AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.roleService.AssignPermissions(c.Request.Context(), id, &req, roleAuthorizationVersion(c))
	result.ResponseResult(c)
}
