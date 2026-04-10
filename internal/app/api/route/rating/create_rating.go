package ratingroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	ratingrequest "manga-go/internal/pkg/request/rating"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Create rating
// @Description  Create a rating for a comic, or update current user's rating if it already exists
// @Tags         Rating
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string                             true  "Comic slug"
// @Param        body  body      ratingrequest.CreateRatingRequest  true  "Rating creation request"
// @Success      200   {object}  response.Result
// @Failure      400   {object}  response.Result
// @Failure      401   {object}  response.Result
// @Failure      404   {object}  response.Result
// @Failure      500   {object}  response.Result
// @Security     AccessToken
// @Router       /ratings/comics/{comicSlug} [post]
func (h *RatingHandler) createRating(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	var req ratingrequest.CreateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}
	comicId, _ := common.GetComicIdFromContext(c.Request.Context())

	result := h.ratingService.CreateRating(c.Request.Context(), user.ID, comicId, &req)
	result.ResponseResult(c)
}
