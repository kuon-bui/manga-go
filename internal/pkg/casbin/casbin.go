package casbin

import (
	"embed"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/logger"
	"strconv"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

//go:embed model.conf
var f embed.FS
var data, _ = f.ReadFile("model.conf")

var modelStr = string(data)

type Enforcer struct {
	*casbin.Enforcer
}

// redisWatcherConfig builds the connection the watcher uses to broadcast and
// receive policy-change notifications. It is a separate connection from the
// application Redis client, so it needs the same credentials.
func redisWatcherConfig(cfg *config.Config) (string, rediswatcher.WatcherOptions) {
	address := cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port)

	return address, rediswatcher.WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Addr:     address,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
		Channel:    "/casbin",
		IgnoreSelf: true,
	}
}

func NewEnforcer(cfg *config.Config, db *gorm.DB, log *logger.Logger) *Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Error("Failed to create Casbin adapter: %v", err)
		panic(err)
	}

	m, err := model.NewModelFromString(modelStr)
	if err != nil {
		log.Error("Failed to create Casbin model: %v", err)
		panic(err)
	}

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		log.Error("Failed to create Casbin enforcer: %v", err)
		panic(err)
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		log.Error("Failed to load Casbin policy: %v", err)
		panic(err)
	}

	address, watcherOptions := redisWatcherConfig(cfg)
	w, err := rediswatcher.NewWatcher(address, watcherOptions)
	if err != nil {
		// Without a watcher every instance keeps serving its own in-memory copy
		// of the policy, so a permission change on one node never reaches the
		// others. Refuse to start rather than run with silently stale policy.
		log.Errorf("Failed to create Casbin watcher: %v", err)
		panic(err)
	}
	if err := enforcer.SetWatcher(w); err != nil {
		log.Error("Failed to set Casbin watcher: %v", err)
		panic(err)
	}
	if err := w.SetUpdateCallback(func(s string) {
		log.Info("Received Casbin policy update notification: %s", s)
		if err := enforcer.LoadPolicy(); err != nil {
			log.Error("Failed to reload Casbin policy: %v", err)
		}
	}); err != nil {
		log.Error("Failed to set Casbin watcher callback: %v", err)
		panic(err)
	}

	return &Enforcer{enforcer}
}
