package permissionroute

import (
	"manga-go/internal/app/api/common/response"
	permissionrequest "manga-go/internal/pkg/request/permission"

	"github.com/gin-gonic/gin"
)

// @Summary      Create permission
// @Description  Create a new permission
// @Tags         Permission
// @Accept       json
// @Produce      json
// @Param        body  body      permissionrequest.CreatePermissionRequest  true  "Permission creation request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /permissions [post]
// @Security     AccessToken
func (h *PermissionHandler) createPermission(c *gin.Context) {
	var req permissionrequest.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.permissionService.CreatePermission(c.Request.Context(), &req)
	result.ResponseResult(c)
}
