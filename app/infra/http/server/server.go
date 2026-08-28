package server

import (
	"context"
	"net/http"
	"time"

	"github.com/anoop-dryad/canopy/app/config"
	"github.com/anoop-dryad/canopy/app/infra/http/handlers"
	"github.com/anoop-dryad/canopy/app/infra/http/routes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	http *http.Server
}

func NewServer(httpConfig config.HTTP, appConfig config.App, appLog *zap.Logger, deps handlers.Dependencies) *Server {
	engine := gin.New()
	engine.Use(gin.Recovery())
	routes.Register(engine, deps, appConfig.IsProduction)
	return &Server{
		http: &http.Server{
			Addr:           httpConfig.Addr,
			Handler:        engine,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1 MiB
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		if err := s.http.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done(): // OS signal received
		return s.http.Shutdown(context.Background())
	}
}
