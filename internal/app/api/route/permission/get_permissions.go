package permissionroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      List permissions
// @Description  Get paginated list of permissions
// @Tags         Permission
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Success      200    {object}  response.PaginationResponse
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Router       /permissions [get]
// @Security     AccessToken
func (h *PermissionHandler) getPermissions(c *gin.Context) {
	var paging common.Paging
	if err := c.ShouldBindQuery(&paging); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.permissionService.ListPermissions(c.Request.Context(), &paging)
	result.ResponseResult(c)
}
