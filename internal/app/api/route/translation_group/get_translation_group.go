package translationgrouproute

import (
	"github.com/gin-gonic/gin"
	_ "manga-go/internal/app/api/common/response"
)

// @Summary      Get translation group by slug
// @Description  Retrieve a specific translation group by its slug
// @Tags         TranslationGroup
// @Accept       json
// @Produce      json
// @Param        slug  path      string  true  "Translation group slug"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Router       /translation-groups/{slug} [get]
// @Security     AccessToken
func (h *TranslationGroupHandler) getTranslationGroup(c *gin.Context) {
	slug := c.Param("translationGroupSlug")

	result := h.translationGroupService.GetTranslationGroup(c.Request.Context(), slug)
	result.ResponseResult(c)
}
