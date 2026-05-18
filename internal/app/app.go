package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func Run(server *http.Server, log *zap.Logger, shutdownFns ...func()) error {
	go func() {
		log.Info("http server started",
			zap.String("addr", server.Addr),
		)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("http server failed",
				zap.Error(err),
			)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stop

	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("http shutdown failed",
			zap.Error(err),
		)

		return err
	}

	for _, fn := range shutdownFns {
		fn()
	}

	log.Info("application stopped")

	return nil
}
