// apps/api/cmd/api/main.go
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"gogym-api/configs"
	gormdb "gogym-api/internal/adapter/db/gorm"

	"github.com/labstack/echo/v4"
)

func main() {
	// ロガー
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Config読込
	cfg, err := configs.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// DB初期化
	_, err = gormdb.NewDB(cfg.Database)
	if err != nil {
		logger.Error("failed to connect DB", "error", err)
		os.Exit(1)
	}
	// TODO: repository層に db を渡す

	// Echoサーバ
	e := echo.New()
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	addr := fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port)
	logger.Info("server starting", "addr", addr, "env", cfg.Server.Env)
	if err := e.Start(addr); err != nil {
		logger.Error("server stopped", "error", err)
	}
}
