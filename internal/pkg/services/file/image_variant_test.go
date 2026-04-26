package fileservice

import "testing"

func TestParseImageVariant(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      ImageVariant
		wantError bool
	}{
		{name: "default to sharp", input: "", want: ImageVariantSharp},
		{name: "economy", input: "economy", want: ImageVariantEconomy},
		{name: "small", input: "small", want: ImageVariantSmall},
		{name: "clear", input: "clear", want: ImageVariantClear},
		{name: "sharp", input: "sharp", want: ImageVariantSharp},
		{name: "invalid", input: "x", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseImageVariant(tt.input)
			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestBuildVariantObjectKey(t *testing.T) {
	canonical := "comics/demo/chapters/ch-1/pages/abc.webp"

	if got := BuildVariantObjectKey(canonical, ImageVariantSharp); got != canonical {
		t.Fatalf("sharp should keep canonical key, got %q", got)
	}

	if got := BuildVariantObjectKey(canonical, ImageVariantEconomy); got != "comics/demo/chapters/ch-1/pages/abc__economy.webp" {
		t.Fatalf("unexpected economy key: %q", got)
	}
}

func TestSanitizeObjectKey(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantError bool
	}{
		{name: "valid", input: "comics/a/cover/abc.webp", want: "comics/a/cover/abc.webp"},
		{name: "leading slash", input: "/comics/a/cover/abc.webp", want: "comics/a/cover/abc.webp"},
		{name: "normalize dot path", input: "comics/a/./cover/abc.webp", want: "comics/a/cover/abc.webp"},
		{name: "empty", input: "", wantError: true},
		{name: "parent traversal", input: "../secret.webp", wantError: true},
		{name: "backslash", input: "comics\\a\\cover.webp", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitizeObjectKey(tt.input)
			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("want %q, got %q", tt.want, got)
			}
		})
	}
}
