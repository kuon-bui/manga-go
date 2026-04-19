package ratingroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      Get user rating
// @Description  Get the current user's rating for the specified comic
// @Tags         Rating
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string  true   "Comic slug"
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

	comicID, ok := common.GetComicIdFromContext(c.Request.Context())
	if !ok {
		response.ResultError("invalid comic id").ResponseResult(c)
		return
	}

	result := h.ratingService.GetUserRatingForComic(c.Request.Context(), user.ID, comicID)
	result.ResponseResult(c)
}
