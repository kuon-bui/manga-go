package userroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      List followed comics
// @Description  Retrieve paginated followed comics for the current user
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Failure      500    {object}  response.Response
// @Security     AccessToken
// @Router       /users/me/followed-comics [get]
func (h *userHandler) getFollowedComics(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	var paging common.Paging
	if err := c.ShouldBindQuery(&paging); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.comicService.ListFollowedComics(c.Request.Context(), user.ID, &paging)
	result.ResponseResult(c)
}
