package ratingroute

import (
	"manga-go/internal/app/api/common/response"
	ratingrequest "manga-go/internal/pkg/request/rating"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Update rating
// @Description  Update current user's rating by ID
// @Tags         Rating
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string                       true  "Comic slug"
// @Param        id    path      string                       true  "Rating ID"
// @Param        body  body      ratingrequest.UpdateRatingRequest  true  "Rating update request"
// @Success      200   {object}  response.Result
// @Failure      400   {object}  response.Result
// @Failure      401   {object}  response.Result
// @Failure      404   {object}  response.Result
// @Failure      500   {object}  response.Result
// @Router       /ratings/comics/{comicSlug}/{id} [put]
// @Security     AccessToken
func (h *RatingHandler) updateRating(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("Invalid id").ResponseResult(c)
		return
	}

	var req ratingrequest.UpdateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.ratingService.UpdateRating(c.Request.Context(), user.ID, id, &req)
	result.ResponseResult(c)
}
