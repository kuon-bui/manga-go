package permissionroute

import (
	_ "manga-go/internal/app/api/common/response"
	"github.com/gin-gonic/gin"
)

// @Summary      Get all permissions
// @Description  Get all permissions without pagination
// @Tags         Permission
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /permissions/all [get]
// @Security     AccessToken
func (h *PermissionHandler) getAllPermissions(c *gin.Context) {
	result := h.permissionService.ListAllPermissions(c.Request.Context())
	result.ResponseResult(c)
}
