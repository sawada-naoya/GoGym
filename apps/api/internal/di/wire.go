//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"log/slog"

	"github.com/google/wire"
	"gorm.io/gorm"

	"gogym-api/configs"
	"gogym-api/internal/adapter/auth"
	"gogym-api/internal/adapter/db/gorm"
	"gogym-api/internal/adapter/http"
	gymUC "gogym-api/internal/usecase/gym"
	userUC "gogym-api/internal/usecase/user"
)

// Application represents the entire application
type Application struct {
	Server *http.Server
	DB     *gorm.DB
	Logger *slog.Logger
}

// Application config is now in configs package

// Cleanup performs application cleanup
func (app *Application) Cleanup() error {
	if sqlDB, err := app.DB.DB(); err == nil {
		sqlDB.Close()
	}
	return nil
}

// ProviderSets - 全てのプロバイダーをここに集約

// InfrastructureSet - インフラ層
var InfrastructureSet = wire.NewSet(
	// Database
	gorm.NewGormDB,
	gorm.NewUserRepository,
	gorm.NewRefreshTokenRepository,
	gorm.NewGymRepository, 
	gorm.NewTagRepository,
	
	// Auth Services
	auth.NewPasswordService,
	auth.NewTokenService,
	
	// S3, Redis, etc. (future)
	// s3.NewS3Service,
	// redis.NewRedisClient,
)

// UseCaseSet - アプリケーション層
var UseCaseSet = wire.NewSet(
	userUC.NewUseCase,
	gymUC.NewUseCase,
)

// HTTPSet - プレゼンテーション層  
var HTTPSet = wire.NewSet(
	http.NewServer,
	http.NewHandler,
	http.NewMiddleware,
)

// Interface Bindings - インターフェースとの紐付け
var InterfaceSet = wire.NewSet(
	// User repository interfaces
	wire.Bind(new(userUC.Repository), new(*gorm.UserRepository)),
	wire.Bind(new(userUC.RefreshTokenRepository), new(*gorm.RefreshTokenRepository)),
	
	// Gym repository interfaces  
	wire.Bind(new(gymUC.Repository), new(*gorm.GymRepository)),
	wire.Bind(new(gymUC.TagRepository), new(*gorm.TagRepository)),
	
	// Auth service interfaces
	wire.Bind(new(userUC.PasswordService), new(*auth.PasswordService)),
	wire.Bind(new(userUC.TokenService), new(*auth.TokenService)),
)

// AllSet - 全プロバイダー
var AllSet = wire.NewSet(
	InfrastructureSet,
	UseCaseSet,
	HTTPSet,
	InterfaceSet,
)

// InitializeApplication - メインのDIコンテナ
func InitializeApplication(ctx context.Context, config *configs.Config, logger *slog.Logger) (*Application, error) {
	wire.Build(
		AllSet,
		
		// Config providers - configs.Configの各フィールドを提供
		wire.FieldsOf(new(*configs.Config), "Database", "Auth", "Server", "S3", "Redis"),
		
		// Application constructor
		wire.Struct(new(Application), "Server", "DB", "Logger"),
	)
	return &Application{}, nil
}