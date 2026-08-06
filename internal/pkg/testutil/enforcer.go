package testutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	casbinpkg "manga-go/internal/pkg/casbin"

	casbinlib "github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
)

// NewInMemoryEnforcer builds an enforcer from the real model.conf with no
// adapter behind it. Tests get the production matcher without a database.
func NewInMemoryEnforcer(t testing.TB) *casbinpkg.Enforcer {
	t.Helper()

	data, err := os.ReadFile(modelConfPath(t))
	if err != nil {
		t.Fatalf("failed to read casbin model: %v", err)
	}

	m, err := model.NewModelFromString(string(data))
	if err != nil {
		t.Fatalf("failed to parse casbin model: %v", err)
	}

	e, err := casbinlib.NewEnforcer(m)
	if err != nil {
		t.Fatalf("failed to create casbin enforcer: %v", err)
	}

	return &casbinpkg.Enforcer{Enforcer: e}
}

// modelConfPath resolves the model relative to this file rather than the test's
// working directory, so callers in any package get the same one.
func modelConfPath(t testing.TB) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to locate the testutil package on disk")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "casbin", "model.conf")
}
