package translationgrouproute

import (
	"manga-go/internal/app/api/common/response"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"
	"manga-go/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h *TranslationGroupHandler) kickMember(c *gin.Context) {
	slug := c.Param("slug")

	var req translationgrouprequest.KickMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	currentUser, err := utils.GetCurrentUserFromGinContext(c)
	if err != nil {
		response.ResultUnauthorized().ResponseResult(c)
		return
	}

	result := h.translationGroupService.KickMember(c.Request.Context(), currentUser.ID, slug, &req)
	result.ResponseResult(c)
}
