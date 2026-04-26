package userroute

import (
	"manga-go/internal/app/api/common/response"
	userrequest "manga-go/internal/pkg/request/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Update user profile
// @Description  Update profile information of a user by id
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id    path      string                           true  "User ID"
// @Param        body  body      userrequest.UpdateUserProfileRequest  true  "User profile update request"
// @Success      200   {object}  response.Result
// @Failure      400   {object}  response.Result
// @Failure      401   {object}  response.Result
// @Failure      500   {object}  response.Result
// @Router       /users/{id} [patch]
// @Security     AccessToken
func (h *userHandler) updateUserProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("Invalid id").ResponseResult(c)
		return
	}

	var req userrequest.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.userService.UpdateUserProfile(c.Request.Context(), userID, &req)
	result.ResponseResult(c)
}
