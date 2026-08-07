package userroute

import (
	"strings"

	"manga-go/internal/app/api/common/response"
	userrequest "manga-go/internal/pkg/request/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Assign roles to user
// @Description  Assign one or more roles to a user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id    path      string                       true  "User ID"
// @Param        body  body      userrequest.AssignRoleRequest  true  "Roles to assign"
// @Param        If-Match  header  string  false  "Current authorization version (for example g12:u4)"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Failure      409   {object}  response.Response  "AUTHORIZATION_STATE_CHANGED, SELF_MANAGE_REQUIRED, or LAST_ROLE_MANAGER"
// @Router       /users/{id}/roles [post]
// @Security     AccessToken
func (h *userHandler) assignUserRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	var req userrequest.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.userService.AssignRoles(c.Request.Context(), id, *req.RoleIDs, authorizationVersion(c))
	result.ResponseResult(c)
}

func authorizationVersion(c *gin.Context) string {
	return strings.Trim(c.GetHeader("If-Match"), "\"")
}
