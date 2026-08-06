package casbin

import (
	"testing"

	"manga-go/internal/pkg/config"
)

func TestRedisWatcherConfigUsesConfiguredAddress(t *testing.T) {
	address, _ := redisWatcherConfig(&config.Config{
		Redis: config.RedisConfig{Host: "redis.internal", Port: 6380},
	})

	if address != "redis.internal:6380" {
		t.Fatalf("expected the watcher to target redis.internal:6380, got %q", address)
	}
}

// The watcher opens its own Redis connection, separate from the application
// client. If it does not carry the configured credentials it silently fails to
// connect and policy changes stop propagating between instances.
func TestRedisWatcherConfigCarriesCredentials(t *testing.T) {
	_, options := redisWatcherConfig(&config.Config{
		Redis: config.RedisConfig{Host: "redis.internal", Port: 6379, Password: "s3cret", DB: 3},
	})

	if options.Options.Password != "s3cret" {
		t.Errorf("expected the configured password to reach the watcher, got %q", options.Options.Password)
	}
	if options.Options.DB != 3 {
		t.Errorf("expected the configured db index to reach the watcher, got %d", options.Options.DB)
	}
}

func TestRedisWatcherConfigIgnoresSelfOriginatedUpdates(t *testing.T) {
	_, options := redisWatcherConfig(&config.Config{
		Redis: config.RedisConfig{Host: "redis.internal", Port: 6379},
	})

	if !options.IgnoreSelf {
		t.Error("expected the watcher to ignore its own updates, otherwise every local write triggers a reload")
	}
	if options.Channel != "/casbin" {
		t.Errorf("expected the watcher to listen on /casbin, got %q", options.Channel)
	}
}
