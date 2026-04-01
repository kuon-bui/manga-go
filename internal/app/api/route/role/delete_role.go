package roleroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Delete role
// @Description  Delete a role by ID
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Role ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /roles/{id} [delete]
// @Security     AccessToken
func (h *RoleHandler) deleteRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	result := h.roleService.DeleteRole(c.Request.Context(), id)
	result.ResponseResult(c)
}
