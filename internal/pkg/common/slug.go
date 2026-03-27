package common

import (
	"context"

	"github.com/google/uuid"
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
