package permissionroute

import (
	_ "manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

// @Summary      List permissions
// @Description  List every permission that can be granted to a role. The catalog is defined in code, so it is fixed for a given deployment.
// @Tags         Permission
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Result
// @Failure      401  {object}  response.Result
// @Failure      403  {object}  response.Result
// @Router       /permissions [get]
// @Security     AccessToken
func (h *PermissionHandler) getPermissions(c *gin.Context) {
	result := h.permissionService.ListPermissions(c.Request.Context())
	result.ResponseResult(c)
}
