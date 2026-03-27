package common

import (
	"context"

	"github.com/google/uuid"
)

const ComicID = "comic-id"

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
