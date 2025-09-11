// di/sets.go
// 役割: 依存性注入のためのwire.NewSetを定義（API Layer）
// Clean Architectureの依存フローに従って各層のプロバイダーを整理
package di

import (
	"github.com/google/wire"

	"gogym-api/configs"
	"gogym-api/internal/adapter/auth"
	gormAdapter "gogym-api/internal/adapter/db/gorm"
	"gogym-api/internal/adapter/http/handler"
	"gogym-api/internal/adapter/http/router"
	"gogym-api/internal/infra/db"
	userUC "gogym-api/internal/usecase/user"
	gymUC "gogym-api/internal/usecase/gym"
)

// =============================================================================
// 1. Infrastructure Set (最下位層)
// Config → Logger, DB, Redis, S3, Auth Services
// =============================================================================

// InfrastructureSet は外部サービスとの接続を提供
// データベース接続、認証サービス、将来のRedis/S3等
var InfrastructureSet = wire.NewSet(
	// 設定から各コンポーネント用の設定を抽出
	ProvideDatabaseConfig,
	
	// データベース接続
	db.NewGormDB,
	
	// 認証サービス（パスワードハッシュ化、JWT生成）
	auth.NewPasswordService,
	auth.NewTokenService,
	
	// インターフェースバインディング（最小限）
	wire.Bind(new(userUC.PasswordService), new(*auth.PasswordService)),
	wire.Bind(new(userUC.TokenService), new(*auth.TokenService)),
	
	// 将来の拡張: Redis, S3等
	// redis.NewRedisClient,
	// s3.NewS3Service,
)

// ProvideDatabaseConfig は設定からデータベース設定を提供
func ProvideDatabaseConfig(cfg *configs.Config) configs.DatabaseConfig {
	return cfg.Database
}

// =============================================================================
// 2. Repository Set (データアクセス層)
// Infrastructure → Repository implementations (GORM/Redis/S3)
// =============================================================================

// RepositorySet はデータアクセス層の実装を提供
// GORM実装、Redis実装、S3実装等
var RepositorySet = wire.NewSet(
	// Repository実装
	gormAdapter.NewUserRepository,
	gormAdapter.NewRefreshTokenRepository,
	gormAdapter.NewGymRepository,
	gormAdapter.NewTagRepository,
	
	// 将来の拡張
	// redisAdapter.NewCacheRepository,
	// s3Adapter.NewFileRepository,
)

// =============================================================================
// 3. UseCase Set (ビジネスロジック層)
// Repository interfaces + Service interfaces → UseCase implementations
// =============================================================================

// UseCaseSet はビジネスロジック層を提供
// ドメインロジックを実装するユースケース群
var UseCaseSet = wire.NewSet(
	// UseCase実装
	userUC.NewUseCase,
	gymUC.NewUseCase,
)

// =============================================================================
// 4. Handler Set (プレゼンテーション層)  
// UseCase → HTTP Handlers
// =============================================================================

// HandlerSet はHTTPハンドラー層を提供
// HTTPリクエストを受け取りUseCaseに処理を委譲
var HandlerSet = wire.NewSet(
	// Handler実装
	handler.NewUserHandler,
	handler.NewGymHandler,
	handler.NewReviewHandler,
	handler.NewFavoriteHandler,
)

// =============================================================================
// 5. Middleware Set (認証・認可層)
// Auth + その他 → Middleware implementations
// =============================================================================

// MiddlewareSet はHTTPミドルウェア層を提供  
//認証、ログ、CORS等のミドルウェア群
var MiddlewareSet = wire.NewSet(
	// TODO: Middleware実装を追加
	// middleware.NewAuthMiddleware,
	// middleware.NewCORSMiddleware,
	// middleware.NewLogMiddleware,
)

// =============================================================================
// 6. Router Set (ルーティング層)
// Handler + Middleware → Router
// =============================================================================

// RouterSet はHTTPルーター層を提供
// エンドポイントとハンドラーのマッピング
var RouterSet = wire.NewSet(
	// Router実装
	router.NewRouter,
)

// =============================================================================
// 7. Server Set (最上位層)
// Router → Echo *Server
// =============================================================================

// ServerSet はHTTPサーバー層を提供
// Echo サーバーの構築と設定
var ServerSet = wire.NewSet(
	// TODO: Server実装を追加
	// server.NewEchoServer,
)

// =============================================================================
// Interface Bindings (インターフェース結合)
// 具体的な実装をインターフェースにバインド
// =============================================================================

// InterfaceSet はインターフェースと実装の結合を定義
// 依存性逆転の原則に従ったインターフェース結合
// Wireでは、関数が直接インターフェースを返すため、明示的なBindは不要
// RepositorySetで定義された関数が自動的に適切なインターフェースを返す
var InterfaceSet = wire.NewSet(
	// Service interfaces (認証サービス)
	wire.Bind(new(userUC.PasswordService), new(*auth.PasswordService)),
	wire.Bind(new(userUC.TokenService), new(*auth.TokenService)),
)

// =============================================================================
// 全体の依存性組み立て
// =============================================================================

// AllSets は全ての依存性プロバイダーを統合
// Clean Architectureの依存フローに従った順序で構築
var AllSets = wire.NewSet(
	InfrastructureSet, // 1. インフラ層
	RepositorySet,     // 2. リポジトリ層  
	UseCaseSet,        // 3. ユースケース層
	HandlerSet,        // 4. ハンドラー層
	MiddlewareSet,     // 5. ミドルウェア層
	RouterSet,         // 6. ルーター層
	ServerSet,         // 7. サーバー層
	InterfaceSet,      // インターフェース結合
)