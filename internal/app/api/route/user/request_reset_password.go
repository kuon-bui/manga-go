package userroute

import (
	"manga-go/internal/app/api/common/response"
	userrequest "manga-go/internal/pkg/request/user"

	"github.com/gin-gonic/gin"
)

// @Summary      Request password reset
// @Description  Send password reset token to user's email
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body  body      userrequest.RequestResetPasswordRequest  true  "Email address"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /users/request-reset-password [post]
func (h *userHandler) requestResetPassword(c *gin.Context) {
	var req userrequest.RequestResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	res := h.userService.RequestResetPassword(c.Request.Context(), req.Email)

	res.ResponseResult(c)
}
