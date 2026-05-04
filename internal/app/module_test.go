package app_test

import (
	"testing"

	"manga-go/internal/app"
	"manga-go/internal/app/api"
)

func TestAppModulesAreRegistered(t *testing.T) {
	if app.Module == nil {
		t.Fatal("expected app.Module to be non-nil")
	}
	if api.Module == nil {
		t.Fatal("expected api.Module to be non-nil")
	}
}
