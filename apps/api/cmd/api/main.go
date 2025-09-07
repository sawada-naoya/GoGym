// cmd/api/main.go
// 役割: アプリケーションのエントリポイント
// 設定読み込み、ログ初期化、依存性注入、サーバー起動を管理
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gogym-api/configs"
)

func main() {
	// コンテキスト作成
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 設定読み込み
	config, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// ログ初期化
	logger := initLogger(config.Server.Env)
	logger.Info("Starting GoGym API Server")

	// 依存性注入とサーバー初期化
	server, err := InitServer(ctx, config, logger)
	if err != nil {
		logger.Error("Failed to initialize server", "error", err)
		os.Exit(1)
	}

	// グレースフルシャットダウンのためのシグナルハンドリング
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Received shutdown signal")
		
		// タイムアウト付きでシャットダウン
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Failed to shutdown server gracefully", "error", err)
		}
		cancel()
	}()

	// サーバー起動
	logger.Info("Server configuration", 
		"port", config.Server.Port,
		"env", config.Server.Env,
		"db_host", config.Database.Host,
	)

	if err := server.Start(); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}

// initLogger はログレベルと出力形式を設定
func initLogger(env string) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// 開発環境では見やすいテキスト形式
	// 本番環境ではJSON形式
	if env == "development" {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}