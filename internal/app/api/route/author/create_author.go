package authorroute

import (
	"manga-go/internal/app/api/common/response"
	authorrequest "manga-go/internal/pkg/request/author"

	"github.com/gin-gonic/gin"
)

// @Summary      Create author
// @Description  Create a new author
// @Tags         Author
// @Accept       json
// @Produce      json
// @Param        body  body      authorrequest.CreateAuthorRequest  true  "Author creation request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Router       /authors [post]
// @Security     AccessToken
func (h *AuthorHandler) createAuthor(c *gin.Context) {
	var req authorrequest.CreateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.authorService.CreateAuthor(c.Request.Context(), &req)
	result.ResponseResult(c)
}
