package roleroute

import (
	"manga-go/internal/app/api/common/response"
	rolerequest "manga-go/internal/pkg/request/role"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Update role
// @Description  Update role information
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id    path      string                    true  "Role ID"
// @Param        body  body      rolerequest.UpdateRoleRequest  true  "Role update request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /roles/{id} [put]
// @Security     AccessToken
func (h *RoleHandler) updateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	var req rolerequest.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.roleService.UpdateRole(c.Request.Context(), id, &req)
	result.ResponseResult(c)
}
