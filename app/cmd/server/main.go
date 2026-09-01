package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/anoop-dryad/canopy/app/config"
	"github.com/anoop-dryad/canopy/app/infra/db"
	"github.com/anoop-dryad/canopy/app/infra/http/handlers"
	"github.com/anoop-dryad/canopy/app/infra/http/server"
	"github.com/anoop-dryad/canopy/app/infra/logger"
	"github.com/anoop-dryad/canopy/app/internal/repository"
	"github.com/anoop-dryad/canopy/app/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	appLog := newLogger(cfg.App)
	defer func() {
		_ = appLog.Sync() // best-effort flush; nothing actionable on shutdown
	}()

	// context cancelled on OS signal — drives graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// ------------------------------- infrastructure ---------------------------------------

	dbPool := db.NewPostgresPool(cfg.DB)

	// ------------------------------- domain ---------------------------------------

	deviceSvc := service.NewService(repository.NewRepository(dbPool), appLog)

	// ------------------------------- http server (blocks) ---------------------------------------

	deps := handlers.Dependencies{
		DeviceHandler: handlers.NewDeviceHandler(deviceSvc),
	}

	srv := server.NewServer(cfg.HTTP, cfg.App, appLog, deps)
	appLog.Info("starting server", zap.String("addr", cfg.HTTP.Addr))
	if err := srv.Start(ctx); err != nil {
		appLog.Fatal("server failed", zap.Error(err))
	}
	appLog.Info("server stopped gracefully")
}

// ------------------------------- helper functions ---------------------------------------

func newLogger(appConfig config.App) *zap.Logger {
	appLog, err := logger.New(appConfig)
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}

	return appLog
}
