package commentroute

import (
	"manga-go/internal/app/api/common/response"
	commentrequest "manga-go/internal/pkg/request/comment"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Report a comment
// @Description  Report an inappropriate comment for moderation review
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Comment ID"
// @Param        request  body      commentrequest.ReportCommentRequest  true  "Report details"
// @Success      200      {object}  response.Result
// @Failure      400      {object}  response.Result
// @Failure      401      {object}  response.Result
// @Failure      404      {object}  response.Result
// @Failure      500      {object}  response.Result
// @Security     AccessToken
// @Router       /comments/{id}/report [post]
func (h *CommentHandler) reportComment(c *gin.Context) {
	var req commentrequest.ReportCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ResultError("Invalid id").ResponseResult(c)
		return
	}

	user, _ := utils.GetCurrentUserFromGinContext(c)
	if user == nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	result := h.commentService.ReportComment(c.Request.Context(), user.ID, id, &req)
	result.ResponseResult(c)
}
