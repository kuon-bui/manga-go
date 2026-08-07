package userroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/authorization"

	"github.com/gin-gonic/gin"
)

// getMyAuthorization godoc
//
//	@Summary		Get the current user's effective authorization profile
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/users/me/authorization [get]
//	@Security		AccessToken
func (h *userHandler) getMyAuthorization(c *gin.Context) {
	viewer := authorization.ViewerFromContext(c.Request.Context())
	if viewer.User == nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	profile, err := h.authorizationAdmin.GetProfile(c.Request.Context(), viewer.User.ID)
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}
	response.ResultSuccess("Authorization profile retrieved successfully", profile).ResponseResult(c)
}
