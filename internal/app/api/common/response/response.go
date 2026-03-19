package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Err     string `json:"error,omitempty"`
}

type PaginationResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Total   int    `json:"total"`
}

type Result struct {
	Success    bool
	HttpStatus int
	Message    string
	Data       any
	Error      error
}

func NewResult(success bool, status int, message string, data any, err error) Result {
	return Result{
		Success:    success,
		HttpStatus: status,
		Message:    message,
		Data:       data,
		Error:      err,
	}
}

func ResponsePaginationData(elements any, total int64) map[string]any {
	return map[string]any{
		"elements": elements,
		"total":    total,
	}
}

func (result Result) ResponseResult(c *gin.Context) {
	// Lấy tracer
	tracer := otel.Tracer("response-result")

	// Lấy context từ request
	_, span := tracer.Start(c.Request.Context(), "ResponseResult")
	defer span.End()

	span.SetAttributes(
		attribute.Bool("success", result.Success),
		attribute.Int("httpStatus", result.HttpStatus),
		attribute.String("message", result.Message),
	)

	if result.Error != nil {
		span.RecordError(result.Error)
		span.SetAttributes(attribute.String("error", result.Error.Error()))
	}

	c.JSON(result.HttpStatus, Response{
		Message: result.Message,
		Data:    result.Data,
		Err: func() string {
			if result.Error != nil {
				return result.Error.Error()
			} else {
				return ""
			}
		}(),
	})
}

func ResponseBadRequest(c *gin.Context, message ...string) {
	c.JSON(http.StatusBadRequest, Response{
		Message: strings.Join(message, ", "),
	})
}

func ResponseNotFound(c *gin.Context, entity string) {
	ResultNotFound(entity).ResponseResult(c)
}

func ResultNotFound(entity string) Result {
	return NewResult(
		false,
		http.StatusBadRequest,
		entity+" not found",
		nil,
		nil,
	)
}

func ResponseUnauthorized(c *gin.Context) {
	ResultUnauthorized().ResponseResult(c)
}

func ResultUnauthorized() Result {
	return NewResult(
		false,
		http.StatusUnauthorized,
		"unauthorized",
		nil,
		nil,
	)
}

func ResultError(message string) Result {
	return NewResult(
		false,
		http.StatusBadRequest,
		message,
		nil,
		nil,
	)
}

func ResultErrDb(err error) Result {
	return NewResult(
		false,
		http.StatusInternalServerError,
		"database error",
		nil,
		err,
	)
}

func ResultErrInternal(err error) Result {
	return NewResult(
		false,
		http.StatusInternalServerError,
		"internal server error",
		nil,
		err,
	)
}

func ResultSuccess(message string, data any) Result {
	return NewResult(
		true,
		http.StatusOK,
		message,
		data,
		nil,
	)
}

func ResultInvalidRequestErr(err error) Result {
	return NewResult(
		false,
		http.StatusBadRequest,
		"invalid request",
		nil,
		err,
	)
}
