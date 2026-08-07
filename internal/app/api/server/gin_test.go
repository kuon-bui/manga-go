package server

import (
	"slices"
	"testing"

	"manga-go/internal/pkg/config"
)

func TestBuildCorsConfigAllowsOptimisticConcurrencyHeader(t *testing.T) {
	t.Parallel()

	corsConfig := buildCorsConfig(&config.Config{})
	if !slices.Contains(corsConfig.AllowHeaders, "If-Match") {
		t.Fatalf("expected CORS allow headers to include If-Match, got %v", corsConfig.AllowHeaders)
	}
}
