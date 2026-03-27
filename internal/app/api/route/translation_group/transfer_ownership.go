package translationgrouproute

import (
	"manga-go/internal/app/api/common/response"
	translationgrouprequest "manga-go/internal/pkg/request/translation_group"

	"github.com/gin-gonic/gin"
)

func (h *TranslationGroupHandler) transferOwnership(c *gin.Context) {
	slug := c.Param("slug")

	var req translationgrouprequest.TransferOwnershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	result := h.translationGroupService.TransferOwnership(c.Request.Context(), slug, &req)
	result.ResponseResult(c)
}
