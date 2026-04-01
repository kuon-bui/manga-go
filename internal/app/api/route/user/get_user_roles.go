package userroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Get user roles
// @Description  Retrieve all roles assigned to a user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "User ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /users/{id}/roles [get]
// @Security     AccessToken
func (h *userHandler) getUserRoles(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	result := h.userService.GetUserRoles(c.Request.Context(), id)
	result.ResponseResult(c)
}
