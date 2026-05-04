package common

import (
	"context"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	ComicID            = "comic-id"
	ChapterID          = "chapter-id"
	TranslationGroupID = "translation-group-id"
	GenreID            = "genre-id"
	TagID              = "tag-id"
)

func setValueToContext(ctx context.Context, key string, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

func getValueFromContext[T any](ctx context.Context, key string) (T, bool) {
	value, ok := ctx.Value(key).(T)
	return value, ok
}

func SetComicIdToContext(ctx context.Context, id uuid.UUID) context.Context {
	return setValueToContext(ctx, ComicID, id)
}

func GetComicIdFromContext(ctx context.Context) (uuid.UUID, bool) {
	return getValueFromContext[uuid.UUID](ctx, ComicID)
}

func SetChapterIdToContext(ctx context.Context, id uuid.UUID) context.Context {
	return setValueToContext(ctx, ChapterID, id)
}

func GetChapterIdFromContext(ctx context.Context) (uuid.UUID, bool) {
	return getValueFromContext[uuid.UUID](ctx, ChapterID)
}

func SetTranslationGroupIdToContext(ctx context.Context, id uuid.UUID) context.Context {
	return setValueToContext(ctx, TranslationGroupID, id)
}

func GetTranslationGroupIdFromContext(ctx context.Context) (uuid.UUID, bool) {
	return getValueFromContext[uuid.UUID](ctx, TranslationGroupID)
}

func SetGenreIdToContext(ctx context.Context, id uuid.UUID) context.Context {
	return setValueToContext(ctx, GenreID, id)
}

func GetGenreIdFromContext(ctx context.Context) (uuid.UUID, bool) {
	return getValueFromContext[uuid.UUID](ctx, GenreID)
}

func SetTagIdToContext(ctx context.Context, id uuid.UUID) context.Context {
	return setValueToContext(ctx, TagID, id)
}

func GetTagIdFromContext(ctx context.Context) (uuid.UUID, bool) {
	return getValueFromContext[uuid.UUID](ctx, TagID)
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
