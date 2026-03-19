package userroute

import (
	"base-go/internal/app/api/common/response"
	userrequest "base-go/internal/pkg/request/user"

	"github.com/gin-gonic/gin"
)

func (h *userHandler) requestResetPassword(c *gin.Context) {
	var req userrequest.RequestResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	res := h.userService.RequestResetPassword(c.Request.Context(), req.Email)

	res.ResponseResult(c)
}
