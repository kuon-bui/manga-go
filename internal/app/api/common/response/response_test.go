package response_test

import (
	"errors"
	"net/http"
	"testing"

	"manga-go/internal/app/api/common/response"

	"github.com/stretchr/testify/assert"
)

func TestResultSuccess(t *testing.T) {
	data := map[string]string{"key": "value"}
	result := response.ResultSuccess("operation successful", data)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "operation successful", result.Message)
	assert.Equal(t, data, result.Data)
	assert.Nil(t, result.Error)
}

func TestResultSuccess_NilData(t *testing.T) {
	result := response.ResultSuccess("deleted", nil)

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "deleted", result.Message)
	assert.Nil(t, result.Data)
}

func TestResultNotFound(t *testing.T) {
	result := response.ResultNotFound("Author")

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "Author not found", result.Message)
	assert.Nil(t, result.Data)
	assert.Nil(t, result.Error)
}

func TestResultUnauthorized(t *testing.T) {
	result := response.ResultUnauthorized()

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusUnauthorized, result.HttpStatus)
	assert.Equal(t, "unauthorized", result.Message)
}

func TestResultError(t *testing.T) {
	result := response.ResultError("something went wrong")

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "something went wrong", result.Message)
}

func TestResultErrDb(t *testing.T) {
	err := errors.New("connection refused")
	result := response.ResultErrDb(err)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	assert.Equal(t, "database error", result.Message)
	assert.Equal(t, err, result.Error)
}

func TestResultErrInternal(t *testing.T) {
	err := errors.New("unexpected panic")
	result := response.ResultErrInternal(err)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusInternalServerError, result.HttpStatus)
	assert.Equal(t, "internal server error", result.Message)
	assert.Equal(t, err, result.Error)
}

func TestResultPaginationData(t *testing.T) {
	items := []string{"a", "b", "c"}
	result := response.ResultPaginationData(items, 10, "items retrieved")

	assert.True(t, result.Success)
	assert.Equal(t, http.StatusOK, result.HttpStatus)
	assert.Equal(t, "items retrieved", result.Message)
	assert.NotNil(t, result.Data)
}

func TestNewResult(t *testing.T) {
	err := errors.New("test error")
	result := response.NewResult(false, http.StatusForbidden, "forbidden", nil, err)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusForbidden, result.HttpStatus)
	assert.Equal(t, "forbidden", result.Message)
	assert.Nil(t, result.Data)
	assert.Equal(t, err, result.Error)
}

func TestResultInvalidRequestErr(t *testing.T) {
	err := errors.New("field is required")
	result := response.ResultInvalidRequestErr(err)

	assert.False(t, result.Success)
	assert.Equal(t, http.StatusBadRequest, result.HttpStatus)
	assert.Equal(t, "invalid request", result.Message)
	assert.Equal(t, err, result.Error)
}
