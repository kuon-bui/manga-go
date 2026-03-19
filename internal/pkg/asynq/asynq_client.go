package asynqclient

import (
	"base-go/internal/pkg/config"
	"fmt"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type asynqClientParams struct {
	fx.In

	Config *config.Config
}

func NewAsynqClient(p asynqClientParams) *asynq.Client {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", p.Config.Redis.Host, p.Config.Redis.Port),
		DB:       p.Config.Redis.DB,
		Password: p.Config.Redis.Password,
		PoolSize: 100,
	})

	return client
}
