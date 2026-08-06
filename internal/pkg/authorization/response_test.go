package authorization

import (
	"context"
	"net/http"
	"testing"
)

func TestEnforceAnyResultDeniesWhenAuthorizerIsMissing(t *testing.T) {
	result := EnforceAnyResult(context.Background(), nil, misconfiguredRequest(), DefaultContexts())

	if result == nil {
		t.Fatal("expected a denying result when the authorizer is missing, got nil (allow)")
	}
	if result.HttpStatus != http.StatusInternalServerError {
		t.Fatalf("expected 500 for a missing authorizer, got %d", result.HttpStatus)
	}
}
