package comicroute

import (
	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

// @Summary      Get trending comics
// @Description  Get list of trending comics based on follow count
// @Tags         Comic
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false  "Number of trending comics (max 50, default 5)"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Result
// @Failure      401    {object}  response.Result
// @Failure      500    {object}  response.Result
// @Security     AccessToken
// @Router       /comics/trending [get]
func (h *ComicHandler) getTrendingComics(c *gin.Context) {
	var req struct {
		Limit int `form:"limit"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.comicService.GetTrendingComics(c.Request.Context(), req.Limit)
	result.ResponseResult(c)
}
