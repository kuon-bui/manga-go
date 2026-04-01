package roleroute

import (
	"manga-go/internal/app/api/common/response"
	rolerequest "manga-go/internal/pkg/request/role"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Assign permissions to role
// @Description  Assign one or more permissions to a role
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id    path      string                          true  "Role ID"
// @Param        body  body      rolerequest.AssignPermissionRequest  true  "Permissions to assign"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
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

	result := h.roleService.AssignPermissions(c.Request.Context(), id, req.PermissionIDs)
	result.ResponseResult(c)
}
