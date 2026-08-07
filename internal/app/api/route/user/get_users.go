package userroute

import (
	"manga-go/internal/app/api/common/response"
	userrequest "manga-go/internal/pkg/request/user"
	authorizationadmin "manga-go/internal/pkg/services/authorization_admin"

	"github.com/gin-gonic/gin"
)

// getUsers godoc
//
//	@Summary		List users for authorization administration
//	@Tags			User
//	@Produce		json
//	@Param			page		query	int		false	"Page"
//	@Param			limit		query	int		false	"Page size"
//	@Param			search		query	string	false	"Name or email"
//	@Param			role_id		query	string	false	"Role ID"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Router			/users [get]
//	@Security		AccessToken
func (h *userHandler) getUsers(c *gin.Context) {
	request := &userrequest.ListAuthorizationUsersRequest{}
	if err := c.ShouldBindQuery(request); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}
	if err := request.Validate(); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	page, err := h.authorizationAdmin.ListUsers(c.Request.Context(), authorizationadmin.ListUsersInput{
		Page:   request.Page,
		Limit:  request.Limit,
		Search: request.Search,
		RoleID: request.RoleID,
	})
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}
	response.ResultPaginationData(page.Data, page.Total, "Users retrieved successfully").ResponseResult(c)
}
