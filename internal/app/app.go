package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func Run() error {
	// --- logger ---
	logg, err := newLogger()
	if err != nil {
		return err
	}

	defer func() {
		_ = logg.Sync()
	}()

	logg.Info("starting application")

	// --- load environment ---
	cfg, err := loadConfig(logg)
	if err != nil {
		return err
	}

	// --- DB ---
	db, err := newPostgres(logg, cfg)
	if err != nil {
		return err
	}

	defer db.Close()

	// --- repositories - services - handlers ---
	deps := initDependencies(logg, db)

	// --- router ---
	router := newRouter(logg, deps)

	// --- HTTP server ---
	srv := newHTTPServer(cfg, router)

	// --- start server ---
	go runServer(logg, srv, cfg.AppPort)

	// --- graceful shutdown ---
	waitShutdownSignal(logg)

	return shutdownServer(logg, srv)
}

func waitShutdownSignal(logg *zap.Logger) {
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	logg.Info("waiting shutdown signal")

	sig := <-quit

	logg.Info(
		"shutdown signal received",
		zap.String("signal", sig.String()),
	)
}

func shutdownServer(logg *zap.Logger, srv *http.Server) error {
	logg.Info("shutdown server")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logg.Error("server shutdown failed",
			zap.Error(err),
		)
		return err
	}

	logg.Info("server stopped")

	return nil
}
