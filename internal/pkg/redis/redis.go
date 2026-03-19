package redis

import (
	"base-go/internal/pkg/logger"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
	logger *logger.Logger
}

var globalRedis *Redis

const Nil = redis.Nil

func NewRedis(client *redis.Client, logger *logger.Logger) *Redis {
	globalRedis = &Redis{client, logger}
	return globalRedis
}

func GetRedis() *Redis {
	if globalRedis == nil {
		panic("Redis client is not initialized")
	}

	return globalRedis
}

func (r *Redis) Client() *redis.Client {
	return r.client
}
