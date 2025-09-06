// GoGym API サーバーの設定管理パッケージ
// 環境変数から設定を読み込み、型安全な設定構造体として提供する
// github.com/caarlos0/env/v10 を使用して環境変数の自動パースを行う
package configs

import (
	"errors"
	"net/url"
	"strings"
	"time"

	env "github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// DatabaseConfig はデータベース接続設定を保持する構造体
// MySQL 8.0への接続に必要な情報を環境変数から取得する
type DatabaseConfig struct {
	Host     string `env:"DB_HOST,required"`        // データベースのホスト名（必須）
	Port     string `env:"DB_PORT" envDefault:"3307"` // データベースのポート番号（デフォルト: 3307）
	User     string `env:"DB_USER,required"`        // データベースユーザー名（必須）
	Password string `env:"DB_PASSWORD,required"`    // データベースパスワード（必須）
	Database string `env:"DB_NAME,required"`        // データベース名（必須）
	Timezone string `env:"TZ" envDefault:"Asia/Tokyo"` // タイムゾーン設定（デフォルト: Asia/Tokyo）
}

// AuthConfig はJWT認証設定を保持する構造体
// アクセストークンとリフレッシュトークンの設定を管理する
type AuthConfig struct {
	JWTSecret        string        `env:"JWT_SECRET,required"`                     // JWTの署名用秘密鍵（必須、16文字以上）
	AccessExpiresIn  time.Duration `env:"JWT_ACCESS_EXPIRES_IN" envDefault:"1h"`   // アクセストークン有効期限（デフォルト: 1時間）
	RefreshExpiresIn time.Duration `env:"JWT_REFRESH_EXPIRES_IN" envDefault:"24h"` // リフレッシュトークン有効期限（デフォルト: 24時間）
	Issuer           string        `env:"JWT_ISSUER" envDefault:"gogym-api"`       // JWTの発行者（デフォルト: gogym-api）
}

// CORSConfig はCross-Origin Resource Sharing設定を保持する構造体
// フロントエンドからのクロスオリジンリクエストを制御する
type CORSConfig struct {
	AllowOrigins []string `env:"CORS_ALLOW_ORIGINS" envSeparator:","`                                          // 許可するオリジン（カンマ区切り）
	AllowMethods []string `env:"CORS_ALLOW_METHODS" envSeparator:"," envDefault:"GET,POST,PUT,DELETE,OPTIONS"` // 許可するHTTPメソッド
	AllowHeaders []string `env:"CORS_ALLOW_HEADERS" envSeparator:"," envDefault:"Content-Type,Authorization"`  // 許可するHTTPヘッダー
	AllowCreds   bool     `env:"CORS_ALLOW_CREDENTIALS" envDefault:"true"`                                     // 認証情報の送信を許可するか
}

// ServerConfig はHTTPサーバー設定を保持する構造体
// Echoサーバーの起動に必要な設定を管理する
type ServerConfig struct {
	Addr string     `env:"APP_ADDR" envDefault:"0.0.0.0"`    // サーバーのバインドアドレス（すべてのインターフェースで待機）
	Port int        `env:"APP_PORT" envDefault:"8081"`       // サーバーのポート番号（デフォルト: 8081）
	Env  string     `env:"APP_ENV" envDefault:"development"` // 環境設定（development/production）
	CORS CORSConfig                                          // CORS設定を埋め込み
}

// S3Config はS3互換オブジェクトストレージ設定を保持する構造体
// MinIOまたはAWS S3への接続に必要な情報を管理する
type S3Config struct {
	Endpoint       *url.URL `env:"S3_ENDPOINT,required"`               // S3エンドポイントURL（必須）
	Bucket         string   `env:"S3_BUCKET,required"`                 // バケット名（必須）
	AccessKey      string   `env:"S3_ACCESS_KEY,required"`             // アクセスキー（必須）
	SecretKey      string   `env:"S3_SECRET_KEY,required"`             // シークレットキー（必須）
	Region         string   `env:"S3_REGION" envDefault:"us-east-1"`   // リージョン（デフォルト: us-east-1）
	PublicURL      *url.URL `env:"S3_PUBLIC_URL"`                      // 公開URL（オプション）
	ForcePathStyle bool     `env:"S3_FORCE_PATH_STYLE" envDefault:"true"` // パススタイルを強制するか（MinIO用）
}

// RedisConfig はRedis接続設定を保持する構造体
// キャッシュとセッション管理のためのRedis接続情報を管理する
type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR" envDefault:"localhost:6380"` // Redis接続アドレス（デフォルト: localhost:6380）
	Password string `env:"REDIS_PASSWORD"`                         // Redis認証パスワード（オプション）
	DB       int    `env:"REDIS_DB" envDefault:"0"`                // 使用するRedisデータベース番号（デフォルト: 0）
}

// Config はアプリケーション全体の設定を保持するメイン構造体
// 各種設定構造体を組み合わせて、アプリケーション起動時に一度だけ読み込まれる
type Config struct {
	Database DatabaseConfig // データベース接続設定
	Auth     AuthConfig     // JWT認証設定
	Server   ServerConfig   // HTTPサーバー設定
	S3       S3Config       // S3互換ストレージ設定
	Redis    RedisConfig    // Redis接続設定
}

// Load は環境変数からアプリケーション設定を読み込む
// .envファイルが存在すれば自動的に読み込み、環境変数をパースして設定構造体を生成する
//
// Returns:
//   - *Config: パース済みの設定構造体
//   - error: パースエラーまたはバリデーションエラー
func Load() (*Config, error) {
	// .envファイルを読み込み（存在しなくてもエラーにしない）
	// 本番環境では環境変数を直接設定するため、.envファイルは主に開発用
	_ = godotenv.Load()

	// 空の設定構造体を作成
	cfg := &Config{}

	// github.com/caarlos0/env/v10 を使用して環境変数を構造体にパース
	// envタグに基づいて自動的に型変換とデフォルト値設定を行う
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// 文字列データの前後空白除去（CORS設定など）
	// 環境変数で設定した値に余計な空白が含まれる場合の対策
	for i, o := range cfg.Server.CORS.AllowOrigins {
		cfg.Server.CORS.AllowOrigins[i] = strings.TrimSpace(o)
	}

	// 設定値のバリデーション実行
	// セキュリティ上重要な設定や、アプリケーションの動作に必要な設定をチェック
	if err := validate(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// validate は読み込まれた設定値の妥当性をチェックする
// セキュリティ要件やアプリケーションの動作要件を満たしているかを検証
//
// Parameters:
//   - c: 検証する設定構造体
//
// Returns:
//   - error: バリデーションエラー（問題なければnil）
func validate(c *Config) error {
	// JWTSecretの長さチェック（セキュリティ要件）
	// 短すぎる秘密鍵は簡単に解読される危険性があるため、最低16文字を要求
	if len(c.Auth.JWTSecret) < 16 {
		return errors.New("JWT_SECRET too short (>=16)")
	}

	// アクセストークン有効期限の妥当性チェック
	// 0以下の値や24時間を超える値は不適切とみなす
	if c.Auth.AccessExpiresIn <= 0 || c.Auth.AccessExpiresIn > 24*time.Hour {
		return errors.New("JWT_ACCESS_EXPIRES_IN out of range (0<ttl<=24h)")
	}

	// S3エンドポイントURLの妥当性チェック
	// nil、スキーム（http/https）なし、ホスト名なしの場合はエラー
	if c.S3.Endpoint == nil || c.S3.Endpoint.Scheme == "" || c.S3.Endpoint.Host == "" {
		return errors.New("S3_ENDPOINT invalid")
	}

	// 本番環境でのCORS設定チェック（セキュリティ要件）
	// 本番環境で全てのオリジンを許可（*）するのはセキュリティリスクが高い
	for _, o := range c.Server.CORS.AllowOrigins {
		if o == "*" && c.Server.Env == "production" {
			return errors.New("CORS_ALLOW_ORIGINS must not be * in production")
		}
	}
	return nil
}
