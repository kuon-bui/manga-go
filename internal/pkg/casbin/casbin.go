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

	if cfg != nil {
		address := cfg.Redis.Host + ":" + strconv.Itoa(cfg.Redis.Port)
		w, err := rediswatcher.NewWatcher(address, rediswatcher.WatcherOptions{
			Options: redis.Options{
				Network: "tcp",
				// Password: os.Getenv("REDIS_PASSWORD"),
			},
			Channel:    "/casbin",
			IgnoreSelf: true,
		})
		if err != nil {
			log.Warnf("Failed to create Casbin Redis watcher: %v", err)
		} else {
			enforcer.SetWatcher(w)
			w.SetUpdateCallback(func(s string) {
				log.Infof("Received Casbin policy update notification: %s", s)
				if err := enforcer.LoadPolicy(); err != nil {
					log.Errorf("Failed to reload Casbin policy: %v", err)
				}
			})
		}
	}

	return &Enforcer{enforcer}
}
