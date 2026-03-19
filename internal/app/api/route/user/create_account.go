package userroute

import (
	"base-go/internal/app/api/common/response"
	userrequest "base-go/internal/pkg/request/user"

	"github.com/gin-gonic/gin"
)

func (h *userHandler) createAccount(c *gin.Context) {
	var req userrequest.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.userService.CreateAccount(c.Request.Context(), &req)
	result.ResponseResult(c)
}
