package comicroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Get comic follow status
// @Description  Retrieve follow status for the current user on a comic
// @Tags         Comic
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string  true  "Comic slug"
// @Success      200        {object}  response.Result
// @Failure      401        {object}  response.Response
// @Failure      404        {object}  response.Response
// @Failure      500        {object}  response.Response
// @Security     AccessToken
// @Router       /comics/{comicSlug}/follow-status [get]
func (h *ComicHandler) getComicFollowStatus(c *gin.Context) {
	user, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	comicID, _ := common.GetComicIdFromContext(c.Request.Context())
	if comicID == uuid.Nil {
		response.ResultNotFound("Comic").ResponseResult(c)
		return
	}

	result := h.comicService.GetFollowStatus(c.Request.Context(), user.ID, comicID)
	result.ResponseResult(c)
}
