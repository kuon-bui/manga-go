package userroute

import (
	"base-go/internal/app/api/common/response"
	userrequest "base-go/internal/pkg/request/user"

	"github.com/gin-gonic/gin"
)

func (h *userHandler) resetPassword(c *gin.Context) {
	var req userrequest.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	res := h.userService.ResetPassword(c.Request.Context(), req)

	res.ResponseResult(c)
}
