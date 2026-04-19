package utils

import (
	"crypto/rand"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func GetRootDir() string {
	workingDir, _ := os.Getwd()

	// loại bỏ phần "/cmd/api" nếu tồn tại
	rootDir := strings.TrimSuffix(workingDir, "/cmd/dev")
	rootDir = strings.TrimSuffix(rootDir, "/cmd/queue")

	return rootDir
}

func TokenGenerator(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// Slugify converts a title to a URL-friendly slug
// Handles Vietnamese characters and special symbols
func Slugify(s string) string {
	// Normalize unicode (decompose Vietnamese characters)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)))
	result, _, _ := transform.String(t, s)

	// Convert to lowercase
	result = strings.ToLower(result)

	// Replace spaces and special chars with hyphens
	result = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`[\s]+`).ReplaceAllString(result, "-")
	result = regexp.MustCompile(`[-]+`).ReplaceAllString(result, "-")

	// Trim hyphens from start and end
	result = strings.Trim(result, "-")

	return result
}
