package roleroute

import (
	_ "manga-go/internal/app/api/common/response"
	"github.com/gin-gonic/gin"
)

// @Summary      Get all roles
// @Description  Get all roles without pagination
// @Tags         Role
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /roles/all [get]
// @Security     AccessToken
func (h *RoleHandler) getAllRoles(c *gin.Context) {
	result := h.roleService.ListAllRoles(c.Request.Context())
	result.ResponseResult(c)
}
