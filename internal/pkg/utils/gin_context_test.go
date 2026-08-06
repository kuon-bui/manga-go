package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"manga-go/internal/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func newTestGinContext() *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return c
}

func TestGetCurrentUserFromGinContextReturnsErrorWhenUnset(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newTestGinContext()

	user, err := GetCurrentUserFromGinContext(c)

	if err == nil {
		t.Fatal("expected an error when no user is set on the context, got nil")
	}
	if user != nil {
		t.Fatalf("expected no user, got %#v", user)
	}
}

func TestGetCurrentUserFromGinContextReturnsErrorForWrongType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newTestGinContext()
	c.Set("current_user", "not-a-user")

	user, err := GetCurrentUserFromGinContext(c)

	if err == nil {
		t.Fatal("expected an error when the context value is not a user, got nil")
	}
	if user != nil {
		t.Fatalf("expected no user, got %#v", user)
	}
}

func TestGetCurrentUserFromGinContextReturnsUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newTestGinContext()
	expected := &model.User{}
	expected.ID = uuid.New()
	SetCurrentUserToGinContext(c, expected)

	user, err := GetCurrentUserFromGinContext(c)

	if err != nil {
		t.Fatalf("expected the stored user to be returned, got error: %v", err)
	}
	if user == nil || user.ID != expected.ID {
		t.Fatalf("expected user %s, got %#v", expected.ID, user)
	}
}
