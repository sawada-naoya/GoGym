//go:build wireinject

// cmd/api/wire.go
// 役割: アプリケーション起動用のDI設定エントリポイント
// wire injectタグでビルド時に依存性注入コードを生成
package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"

	"gogym-api/configs"
)

// Server はEchoサーバーのラッパー構造体
// アプリケーション全体のライフサイクルを管理
type Server struct {
	Echo   *echo.Echo
	Config *configs.Config
	Logger *slog.Logger
}

// InitServer は全ての依存関係を注入してServerを構築
// 
// 依存フロー:
// Config (env) 
//   ├─ Logger(slog)
//   ├─ DB (GORM MySQL) ──┐
//   ├─ Redis              │
//   ├─ S3(MinIO)          │  
//   └─ Auth(JWT)          │
//                         │
//                 Repositories (adapter/db/*) ← DB/Redis/S3に依存
//                         │
//                    Usecases (usecase/*) ← Repo + Auth + その他サービス
//                         │
//                 Handlers (adapter/http/handler/*)
//                         │
//               Middlewares (adapter/http/middleware/*) ← Authなど
//                         │
//                   Router (adapter/http/router.go)
//                         │
//                    Echo *Server
func InitServer(ctx context.Context, config *configs.Config, logger *slog.Logger) (*Server, error) {
	wire.Build(
		// 基本的なEchoサーバーを作成する関数
		NewBasicEcho,
		
		// Server構造体を構築
		wire.Struct(new(Server), "Echo", "Config", "Logger"),
	)
	return &Server{}, nil
}

// NewBasicEcho は基本的なEchoサーバーを作成
func NewBasicEcho() *echo.Echo {
	e := echo.New()
	// 基本的なミドルウェア設定は後で追加
	return e
}

// Start はサーバーを開始
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.Config.Server.Port)
	s.Logger.Info("Starting server", "address", addr)
	return s.Echo.Start(addr)
}

// Shutdown はサーバーを graceful shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info("Shutting down server")
	return s.Echo.Shutdown(ctx)
}