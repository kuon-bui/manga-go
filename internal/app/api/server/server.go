package server

import (
	"errors"
	"fmt"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type NewHttpServerParams struct {
	fx.In

	Config *config.Config
	Logger *logger.Logger
	Router *gin.Engine
}

func NewHttpServer(p NewHttpServerParams) *http.Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", p.Config.Service.Port),
		Handler: p.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			p.Logger.Infof("listen: %s\n", err)
		}

		p.Logger.Info("Server closed")
	}()

	return srv
}
