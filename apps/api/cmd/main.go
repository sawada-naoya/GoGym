package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gogym-api/internal/configs"
	"gogym-api/internal/di"
)

func main() {
	config, err := configs.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
	}

	e, _, err := di.BuildServer(config)
	if err != nil {
		slog.Error("Failed to build server", "error", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", config.HTTP.Host, config.HTTP.Port)
	slog.Info("Starting server", "address", addr)

	errCh := make(chan error, 1)
	go func() {
		errCh <- e.Start(addr)
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-sigCtx.Done():
		slog.Info("shutdown signal received")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			slog.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}
		slog.Info("server shutdown complete")
	case err := <-errCh:
		// 起動直後にエラーで落ちた場合
		if err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}
}
