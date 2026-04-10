package ratingroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      List ratings
// @Description  Get paginated ratings in the current comic context
// @Tags         Rating
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string  true   "Comic slug"
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Result
// @Failure      401    {object}  response.Result
// @Failure      404    {object}  response.Result
// @Failure      500    {object}  response.Result
// @Security     AccessToken
// @Router       /ratings/comics/{comicSlug} [get]
func (h *RatingHandler) getRatings(c *gin.Context) {
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

	result := h.ratingService.ListRatings(c.Request.Context(), user.ID, &paging)
	result.ResponseResult(c)
}
