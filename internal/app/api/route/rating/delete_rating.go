package ratingroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Delete rating
// @Description  Delete current user's rating by ID
// @Tags         Rating
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string  true  "Comic slug"
// @Param        id  path      string  true  "Rating ID"
// @Success      200  {object}  response.Result
// @Failure      400  {object}  response.Result
// @Failure      401  {object}  response.Result
// @Failure      404  {object}  response.Result
// @Failure      500  {object}  response.Result
// @Router       /ratings/comics/{comicSlug}/{id} [delete]
// @Security     AccessToken
func (h *RatingHandler) deleteRating(c *gin.Context) {
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

	result := h.ratingService.DeleteRating(c.Request.Context(), user.ID, id)
	result.ResponseResult(c)
}
