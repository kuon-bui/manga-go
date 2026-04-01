package translationgrouproute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      List translation groups
// @Description  Get paginated list of translation groups
// @Tags         TranslationGroup
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Items per page"
// @Success      200    {object}  response.PaginationResponse
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Router       /translation-groups [get]
// @Security     AccessToken
func (h *TranslationGroupHandler) getTranslationGroups(c *gin.Context) {
	var paging common.Paging
	if err := c.ShouldBindQuery(&paging); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.translationGroupService.ListTranslationGroups(c.Request.Context(), &paging)
	result.ResponseResult(c)
}
