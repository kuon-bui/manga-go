package userroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
