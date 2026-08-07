package userroute

import (
	"errors"

	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// getAuthorizationUser godoc
//
//	@Summary		Get a user authorization snapshot
//	@Tags			User
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Router			/users/{id}/authorization [get]
//	@Security		AccessToken
func (h *userHandler) getAuthorizationUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("invalid id").ResponseResult(c)
		return
	}

	user, err := h.authorizationAdmin.GetUser(c.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.ResultNotFound("User").ResponseResult(c)
		return
	}
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}
	response.ResultSuccess("User authorization retrieved successfully", user).ResponseResult(c)
}
