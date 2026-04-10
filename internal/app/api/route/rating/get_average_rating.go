package ratingroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      Get average rating
// @Description  Get average score and total rating count for a comic
// @Tags         Rating
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string  true  "Comic slug"
// @Success      200        {object}  response.Result
// @Failure      401        {object}  response.Result
// @Failure      404        {object}  response.Result
// @Failure      500        {object}  response.Result
// @Security     AccessToken
// @Router       /ratings/comics/{comicSlug}/average [get]
func (h *RatingHandler) getAverageRating(c *gin.Context) {
	comicId, _ := common.GetComicIdFromContext(c.Request.Context())

	var result response.Result
	result = h.ratingService.GetAverageRating(c.Request.Context(), comicId)
	result.ResponseResult(c)
}
