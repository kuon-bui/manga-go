package permissionroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Delete permission
// @Description  Delete a permission by ID
// @Tags         Permission
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Permission ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Router       /permissions/{id} [delete]
// @Security     AccessToken
func (h *PermissionHandler) deletePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	result := h.permissionService.DeletePermission(c.Request.Context(), id)
	result.ResponseResult(c)
}
