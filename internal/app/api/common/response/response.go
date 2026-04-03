package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func NewResult(success bool, status int, message string, data any, err error) Result {
	return Result{
		Success:    success,
		HttpStatus: status,
		Message:    message,
		Data:       data,
		Error:      err,
	}
}

func responsePaginationData(elements any, total int64) paginationResponse {
	return paginationResponse{
		Data:  elements,
		Total: total,
	}
}

func ResultPaginationData(elements any, total int64, message string) Result {
	return NewResult(
		true,
		http.StatusOK,
		message,
		responsePaginationData(elements, total),
		nil,
	)
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
	errorDetails := ""
	if result.Error != nil {
		errorDetails = result.Error.Error()
	}
	c.JSON(result.HttpStatus, Response{
		Message:          result.Message,
		Data:             result.Data,
		Err:              errorDetails,
		ValidationErrors: result.ValidationErrors,
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
	validationErrors := parseValidationErrors(err)

	return NewResult(
		false,
		http.StatusBadRequest,
		"invalid request",
		nil,
		err,
	).withValidationErrors(validationErrors)
}

func (r Result) withValidationErrors(validationErrors []ValidationFieldError) Result {
	r.ValidationErrors = validationErrors

	return r
}
