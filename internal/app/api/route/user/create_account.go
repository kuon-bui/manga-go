package userroute

import (
	"manga-go/internal/app/api/common/response"
	userrequest "manga-go/internal/pkg/request/user"

	"github.com/gin-gonic/gin"
)

// @Summary      Create new user account
// @Description  Register a new user account with email and password
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body  body      userrequest.CreateUserRequest  true  "User creation request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /users [post]
func (h *userHandler) createAccount(c *gin.Context) {
	var req userrequest.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.userService.CreateAccount(c.Request.Context(), &req)
	result.ResponseResult(c)
}
