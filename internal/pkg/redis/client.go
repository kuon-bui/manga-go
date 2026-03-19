package redis

import (
	"base-go/internal/pkg/config"
	"base-go/internal/pkg/logger"
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(config *config.Config, logger *logger.Logger) *redis.Client {
	cfg := config.Redis
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: 100,
	})

	redisCollector := NewRedisCollector(client)
	prometheus.MustRegister(redisCollector)

	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("Connected to Redis successfully")
	return client
}

type redisCollector struct {
	prometheus.Collector
	rds        *redis.Client
	hit        prometheus.Gauge
	misses     prometheus.Gauge
	timeouts   prometheus.Gauge
	totalConns prometheus.Gauge
	idleConns  prometheus.Gauge
	staleConns prometheus.Gauge
}

func NewRedisCollector(rds *redis.Client) *redisCollector {
	return &redisCollector{
		rds: rds,
		hit: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "redis_hit",
			Help: "Number of times free connection was found in the pool",
		}),
		misses: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "redis_misses",
			Help: "Number of times free connection was not found in the pool",
		}),
		timeouts: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "redis_timeouts",
			Help: "Number of times a wait timeout occurred",
		}),
		totalConns: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "redis_total_conns",
			Help: "Total number of connections in the pool",
		}),
		idleConns: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "redis_idle_conns",
			Help: "Number of idle connections in the pool",
		}),
		staleConns: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "redis_stale_conns",
			Help: "Number of stale connections removed from the pool",
		}),
	}
}

func (c *redisCollector) Describe(ch chan<- *prometheus.Desc) {
	c.hit.Describe(ch)
	c.misses.Describe(ch)
	c.timeouts.Describe(ch)
	c.totalConns.Describe(ch)
	c.idleConns.Describe(ch)
	c.staleConns.Describe(ch)
}

func (c *redisCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.rds.PoolStats()
	c.hit.Set(float64(stats.Hits))
	c.misses.Set(float64(stats.Misses))
	c.timeouts.Set(float64(stats.Timeouts))
	c.totalConns.Set(float64(stats.TotalConns))
	c.idleConns.Set(float64(stats.IdleConns))
	c.staleConns.Set(float64(stats.StaleConns))

	c.hit.Collect(ch)
	c.misses.Collect(ch)
	c.timeouts.Collect(ch)
	c.totalConns.Collect(ch)
	c.idleConns.Collect(ch)
	c.staleConns.Collect(ch)
}
