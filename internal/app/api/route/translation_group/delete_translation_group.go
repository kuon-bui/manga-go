package translationgrouproute

import (
	"github.com/gin-gonic/gin"
	_ "manga-go/internal/app/api/common/response"
)

// @Summary      Delete translation group
// @Description  Delete a translation group by slug
// @Tags         TranslationGroup
// @Accept       json
// @Produce      json
// @Param        slug  path      string  true  "Translation group slug"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /translation-groups/{slug} [delete]
// @Security     AccessToken
func (h *TranslationGroupHandler) deleteTranslationGroup(c *gin.Context) {
	slug := c.Param("translationGroupSlug")

	result := h.translationGroupService.DeleteTranslationGroup(c.Request.Context(), slug)
	result.ResponseResult(c)
}
