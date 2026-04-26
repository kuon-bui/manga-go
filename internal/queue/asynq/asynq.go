package asynq

import (
	"context"
	"fmt"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/logger"
	queueconstant "manga-go/internal/queue/queue_constant"
	"strings"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

type RunAsynqParams struct {
	fx.In

	Config *config.Config
	Logger *logger.Logger
}

func NewAsynqServerMux() *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.Use(tracerMiddleware)

	return mux
}

func NewAsynqServer(p RunAsynqParams) *asynq.Server {
	queues := map[string]int{
		queueconstant.MAIL_DELIVER_QUEUE:  3,
		queueconstant.NOTIFICATION_QUEUE:  5,
		queueconstant.IMAGE_PROCESS_QUEUE: 4,
	}

	var redisConnOpt asynq.RedisConnOpt = asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", p.Config.Redis.Host, p.Config.Redis.Port),
		DB:       p.Config.Redis.DB,
		Password: p.Config.Redis.Password,
	}

	if len(p.Config.Redis.Cluster) > 0 {
		redisConnOpt = asynq.RedisClusterClientOpt{
			Addrs:    strings.Split(p.Config.Redis.Cluster, ","),
			Password: p.Config.Redis.Password,
		}
	}

	server := asynq.NewServer(
		redisConnOpt,
		asynq.Config{
			Concurrency: p.Config.Asynq.Concurrency,
			Queues:      queues,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				p.Logger.Error("asynq error: ", err)
			}),
		},
	)

	return server
}
