package commentroute

import (
	"manga-go/internal/app/api/common/response"
	commentrequest "manga-go/internal/pkg/request/comment"

	"github.com/gin-gonic/gin"
)

// @Summary      List comments on chapter
// @Description  Retrieve paginated comments for a specific chapter
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        chapterId  query  string  true   "Chapter ID"
// @Param        page       query  int     false  "Page number (default: 1)"
// @Param        limit      query  int     false  "Records per page (default: 10)"
// @Success      200        {object}  response.PaginationResponse
// @Failure      400        {object}  response.Response
// @Failure      401        {object}  response.Response
// @Failure      500        {object}  response.Response
// @Security     AccessToken
// @Router       /comments [get]
func (h *CommentHandler) getComments(c *gin.Context) {
	var req commentrequest.ListCommentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.commentService.ListComments(c.Request.Context(), &req)
	result.ResponseResult(c)
}
