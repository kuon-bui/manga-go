package translationgrouproute

import (
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"

	"github.com/gin-gonic/gin"
)

// @Summary      Get translation group members
// @Description  Get a list of members in the translation group
// @Tags         Translation Group
// @Accept       json
// @Produce      json
// @Param        translationGroupSlug path string true "Translation Group slug"
// @Success      200    {object}  response.Result
// @Failure      400    {object}  response.Result
// @Failure      401    {object}  response.Result
// @Failure      500    {object}  response.Result
// @Security     AccessToken
// @Router       /translation-groups/{translationGroupSlug}/members [get]
func (h *TranslationGroupHandler) getMembers(c *gin.Context) {
	id, ok := common.GetTranslationGroupIdFromContext(c.Request.Context())
	if !ok {
		response.ResultError("invalid translation group id").ResponseResult(c)
		return
	}

	result := h.translationGroupService.GetMembers(c.Request.Context(), id)
	result.ResponseResult(c)
}
