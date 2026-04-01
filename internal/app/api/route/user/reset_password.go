package userroute

import (
	"manga-go/internal/app/api/common/response"
	userrequest "manga-go/internal/pkg/request/user"

	"github.com/gin-gonic/gin"
)

// @Summary      Reset user password
// @Description  Reset password using reset token sent to email
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body  body      userrequest.ResetPasswordRequest  true  "Reset password request with token"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /users/reset-password [post]
func (h *userHandler) resetPassword(c *gin.Context) {
	var req userrequest.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	res := h.userService.ResetPassword(c.Request.Context(), req)

	res.ResponseResult(c)
}
