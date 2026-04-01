package roleroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Get role by ID
// @Description  Retrieve a specific role by its ID
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Role ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /roles/{id} [get]
// @Security     AccessToken
func (h *RoleHandler) getRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	result := h.roleService.GetRole(c.Request.Context(), id)
	result.ResponseResult(c)
}
