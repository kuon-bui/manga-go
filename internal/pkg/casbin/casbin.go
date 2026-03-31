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

func NewEnforcer(cfg *config.Config, db *gorm.DB, logger *logger.Logger) *Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		logger.Error("Failed to create Casbin adapter: %v", err)
		panic(err)
	}

	m, _ := model.NewModelFromString(modelStr)
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		logger.Error("Failed to create Casbin enforcer: %v", err)
		panic(err)
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		logger.Error("Failed to load Casbin policy: %v", err)
		panic(err)
	}

	address := cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port)
	w, _ := rediswatcher.NewWatcher(address, rediswatcher.WatcherOptions{
		Options: redis.Options{
			Network: "tcp",
			// Password: os.Getenv("REDIS_PASSWORD"),
		},
		Channel:    "/casbin",
		IgnoreSelf: true,
	})
	enforcer.SetWatcher(w)
	w.SetUpdateCallback(func(s string) {
		logger.Info("Received Casbin policy update notification: %s", s)
		enforcer.LoadPolicy()
	})

	return &Enforcer{enforcer}
}
