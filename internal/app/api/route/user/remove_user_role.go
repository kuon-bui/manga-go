package userroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Remove role from user
// @Description  Remove a role from a user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id      path      string  true  "User ID"
// @Param        roleId  path      string  true  "Role ID"
// @Success      200     {object}  response.Response
// @Failure      400     {object}  response.Response
// @Failure      401     {object}  response.Response
// @Router       /users/{id}/roles/{roleId} [delete]
// @Security     AccessToken
func (h *userHandler) removeUserRole(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid user id").ResponseResult(c)
		return
	}

	roleID, err := uuid.Parse(c.Param("roleId"))
	if err != nil {
		response.ResultError("invalid role id").ResponseResult(c)
		return
	}

	result := h.userService.RemoveRole(c.Request.Context(), userID, roleID)
	result.ResponseResult(c)
}
