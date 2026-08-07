package authorizationroute

import (
	"manga-go/internal/app/api/common/response"
	authorizationrequest "manga-go/internal/pkg/request/authorization"
	authorizationadmin "manga-go/internal/pkg/services/authorization_admin"

	"github.com/gin-gonic/gin"
)

// getAuditLogs godoc
//
//	@Summary		List authorization audit logs
//	@Tags			Authorization
//	@Produce		json
//	@Param			page			query	int		false	"Page"
//	@Param			limit			query	int		false	"Page size"
//	@Param			actor			query	string	false	"Actor name or email"
//	@Param			action			query	string	false	"Action"
//	@Param			target_type	query	string	false	"Target type"
//	@Param			target_id		query	string	false	"Target ID"
//	@Param			start_at		query	string	false	"Start time"
//	@Param			end_at			query	string	false	"End time"
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Router			/authorization/audit-logs [get]
func (h *Handler) getAuditLogs(c *gin.Context) {
	request := &authorizationrequest.ListAuditLogsRequest{}
	if err := c.ShouldBindQuery(request); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}
	if err := request.Validate(); err != nil {
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	page, err := h.authorizationAdmin.ListAuditLogs(c.Request.Context(), authorizationadmin.ListAuditInput{
		Page:       request.Page,
		Limit:      request.Limit,
		Actor:      request.Actor,
		Action:     request.Action,
		TargetType: request.TargetType,
		TargetID:   request.TargetID,
		StartAt:    request.StartAt,
		EndAt:      request.EndAt,
	})
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}
	response.ResultSuccess("Authorization audit logs retrieved successfully", page).ResponseResult(c)
}
