//go:build wireinject

// cmd/api/wire.go
// 役割: アプリケーション起動用のDI設定エントリポイント
// wire injectタグでビルド時に依存性注入コードを生成
package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"gogym-api/configs"
	"gogym-api/internal/adapter/auth"
	gormAdapter "gogym-api/internal/adapter/db/gorm"
	"gogym-api/internal/adapter/http/handler"
	"gogym-api/internal/adapter/http/router"
	"gogym-api/internal/adapter/service"
	"gogym-api/internal/di"
	"gogym-api/internal/infra/db"
	gymUC "gogym-api/internal/usecase/gym"
	reviewUC "gogym-api/internal/usecase/review"
	sessionUC "gogym-api/internal/usecase/session"
	userUC "gogym-api/internal/usecase/user"
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
		// Database
		di.ProvideDatabaseConfig,
		db.NewGormDB,

		// Auth services
		NewJWTService,
		NewBcryptPasswordHasher,
		NewUserPasswordHasher,
		NewUserIDProvider,
		NewSessionIDProvider,
		NewUserRepository,
		NewSessionUserRepository,
		service.NewTimeProvider,

		// Repositories
		gormAdapter.NewGymRepository,
		gormAdapter.NewTagRepository,
		gormAdapter.NewReviewRepository,

		// Use cases
		userUC.NewInteractor,
		gymUC.NewUseCase,
		reviewUC.NewUseCase,
		sessionUC.NewInteractor,

		// Handlers
		handler.NewUserHandler,
		handler.NewGymHandler,
		handler.NewReviewHandler,
		handler.NewFavoriteHandler,
		handler.NewSessionHandler,

		// Router & Server
		NewConfiguredEcho,

		// Server構造体を構築
		wire.Struct(new(Server), "Echo", "Config", "Logger"),
	)
	return &Server{}, nil
}

// NewBasicEcho は基本的なEchoサーバーを作成
func NewBasicEcho() *echo.Echo {
	e := echo.New()
	
	// バリデーターを設定
	e.Validator = &CustomValidator{validator: validator.New()}
	
	return e
}

// CustomValidator カスタムバリデーター
type CustomValidator struct {
	validator *validator.Validate
}

// Validate バリデーション実行
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// NewConfiguredEcho はルーティング設定済みのEchoサーバーを作成
func NewConfiguredEcho(gymHandler *handler.GymHandler, userHandler *handler.UserHandler, reviewHandler *handler.ReviewHandler, favoriteHandler *handler.FavoriteHandler, sessionHandler *handler.SessionHandler) *echo.Echo {
	e := NewBasicEcho()

	// ルーティング設定
	router.NewRouter(e, gymHandler, userHandler, reviewHandler, favoriteHandler, sessionHandler)

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

// NewJWTService はJWTサービスを作成
func NewJWTService(cfg *configs.Config) sessionUC.JWT {
	return auth.New(
		[]byte(cfg.Auth.JWTSecret),
		[]byte(cfg.Auth.JWTSecret),
		cfg.Auth.Issuer,
		cfg.Auth.AccessExpiresIn,
		cfg.Auth.RefreshExpiresIn,
	)
}

// NewBcryptPasswordHasher はBcryptパスワードハッシャーを作成
func NewBcryptPasswordHasher() (sessionUC.PasswordHasher, error) {
	hasher, err := auth.NewBcryptPasswordHasher(12, "")
	if err != nil {
		return nil, err
	}
	return hasher, nil
}

// NewUserPasswordHasher はユーザー用パスワードハッシャーを作成
func NewUserPasswordHasher() (userUC.PasswordHasher, error) {
	hasher, err := auth.NewBcryptPasswordHasher(12, "")
	if err != nil {
		return nil, err
	}
	return hasher, nil
}

// NewUserIDProvider はユーザー用IDプロバイダーを作成
func NewUserIDProvider() userUC.IDProvider {
	return service.NewIDProvider()
}

// NewSessionIDProvider はセッション用IDプロバイダーを作成
func NewSessionIDProvider() sessionUC.IDProvider {
	return service.NewIDProvider()
}

// NewUserRepository はユーザー用ユーザーリポジトリを作成
func NewUserRepository(db *gorm.DB) userUC.Repository {
	return gormAdapter.NewUserRepository(db)
}

// NewSessionUserRepository はセッション用ユーザーリポジトリを作成
func NewSessionUserRepository(db *gorm.DB) sessionUC.UserRepository {
	repo := gormAdapter.NewUserRepository(db)
	return repo.(sessionUC.UserRepository)
}