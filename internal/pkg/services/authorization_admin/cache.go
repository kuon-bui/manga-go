package authorizationadmin

import (
	"context"
	"encoding/json"
	"time"

	rediswrapper "manga-go/internal/pkg/redis"

	redislib "github.com/redis/go-redis/v9"
)

type RedisProfileCache struct {
	client *redislib.Client
}

func NewRedisProfileCache(redisClient *rediswrapper.Redis) *RedisProfileCache {
	return &RedisProfileCache{client: redisClient.Client()}
}

func (c *RedisProfileCache) Get(ctx context.Context, key string, target *AuthorizationProfile) (bool, error) {
	value, err := c.client.Get(ctx, key).Bytes()
	if err == redislib.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(value, target); err != nil {
		return false, err
	}
	return true, nil
}

func (c *RedisProfileCache) Set(ctx context.Context, key string, profile *AuthorizationProfile, ttl time.Duration) error {
	value, err := json.Marshal(profile)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisProfileCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
