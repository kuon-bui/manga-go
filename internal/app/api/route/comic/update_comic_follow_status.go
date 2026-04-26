package comicroute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	comicrequest "manga-go/internal/pkg/request/comic"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Update comic follow status
// @Description  Update follow status for a followed comic of current user
// @Tags         Comic
// @Accept       json
// @Produce      json
// @Param        comicSlug  path      string                     true  "Comic slug"
// @Param        body       body      comicrequest.FollowComicRequest  true  "Follow status update request"
// @Success      200        {object}  response.Result
// @Failure      400        {object}  response.Response
// @Failure      401        {object}  response.Response
// @Failure      404        {object}  response.Response
// @Failure      500        {object}  response.Response
// @Security     AccessToken
// @Router       /comics/{comicSlug}/follow-status [patch]
func (h *ComicHandler) updateComicFollowStatus(c *gin.Context) {
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

	var req comicrequest.FollowComicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.comicService.UpdateComicFollowStatus(c.Request.Context(), user.ID, comicID, req.FollowStatus)
	result.ResponseResult(c)
}
