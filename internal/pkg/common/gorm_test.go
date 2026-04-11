package common

import "testing"

func TestGenerateModelCode(t *testing.T) {
	got := GenerateModelCode(12, 4, "USR")
	if got != "USR0012" {
		t.Fatalf("expected USR0012, got %s", got)
	}
}
