package comicroute

import (
	"manga-go/internal/app/api/common/response"
	comicrequest "manga-go/internal/pkg/request/comic"

	"github.com/gin-gonic/gin"
)

// @Summary      List comics
// @Description  Get paginated list of comics
// @Tags         Comic
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Param        translationGroupSlug  query     string  false  "Translation group slug"
// @Param        genreSlugs  query     []string  false  "Genre slugs"
// @Param        tagSlugs  query     []string  false  "Tag slugs"
// @Param        search  query     string  false  "Search by comic title or alternative titles"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Router       /comics [get]
// @Security     AccessToken
func (h *ComicHandler) getComics(c *gin.Context) {
	var req comicrequest.ListComicsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.comicService.ListComics(c.Request.Context(), &req)
	result.ResponseResult(c)
}
