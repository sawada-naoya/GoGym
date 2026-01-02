package configs

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	env "github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string `env:"DB_HOST,required"`           // データベースのホスト名（必須）
	Port     string `env:"DB_PORT" envDefault:"3307"`  // データベースのポート番号（デフォルト: 3307）
	User     string `env:"DB_USER,required"`           // データベースユーザー名（必須）
	Password string `env:"DB_PASSWORD,required"`       // データベースパスワード（必須）
	Database string `env:"DB_NAME,required"`           // データベース名（必須）
	Timezone string `env:"TZ" envDefault:"Asia/Tokyo"` // タイムゾーン設定（デフォルト: Asia/Tokyo）
}

type AuthConfig struct {
	JWTSecret        string        `env:"JWT_SECRET,required"`                     // JWTの署名用秘密鍵（必須、16文字以上）
	AccessExpiresIn  time.Duration `env:"JWT_ACCESS_EXPIRES_IN" envDefault:"1h"`   // アクセストークン有効期限（デフォルト: 1時間）
	RefreshExpiresIn time.Duration `env:"JWT_REFRESH_EXPIRES_IN" envDefault:"24h"` // リフレッシュトークン有効期限（デフォルト: 24時間）
	Issuer           string        `env:"JWT_ISSUER" envDefault:"gogym-api"`       // JWTの発行者（デフォルト: gogym-api）
}

type CORSConfig struct {
	AllowOrigins []string `env:"CORS_ALLOW_ORIGINS"   envSeparator:","`
	AllowMethods []string `env:"CORS_ALLOW_METHODS"   envSeparator:"," envDefault:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowHeaders []string `env:"CORS_ALLOW_HEADERS"   envSeparator:"," envDefault:"Content-Type,Authorization"`
	AllowCreds   bool     `env:"CORS_ALLOW_CREDENTIALS" envDefault:"true"`
}

type HTTPConfig struct {
	Host string     `env:"APP_HOST" envDefault:"0.0.0.0"` // バインド先
	Port int        `env:"APP_PORT" envDefault:"8081"`    // ポート
	Env  string     `env:"APP_ENV"  envDefault:"development"`
	CORS CORSConfig // 既存のCORS設定をそのまま利用

	// タイムアウト類（http.Server直結）
	ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT"        envDefault:"10s"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`
	WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT"       envDefault:"15s"`
	IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT"        envDefault:"60s"`
	MaxHeaderBytes    int           `env:"HTTP_MAX_HEADER_BYTES"    envDefault:"1048576"` // 1<<20
}

type SlackConfig struct {
	ContactWebhookURL string `env:"SLACK_CONTACT_WEBHOOK_URL"` // 問い合わせ通知用WebhookURL
}

type Config struct {
	Database DatabaseConfig // データベース接続設定
	Auth     AuthConfig     // JWT認証設定
	HTTP     HTTPConfig     // HTTPサーバー設定
	Slack    SlackConfig    // Slack通知設定
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// 前後空白の除去（CORS）
	for i, o := range cfg.HTTP.CORS.AllowOrigins {
		cfg.HTTP.CORS.AllowOrigins[i] = strings.TrimSpace(o)
	}

	// Render互換: PORT を優先（APP_PORTより上位）
	if v := strings.TrimSpace(os.Getenv("PORT")); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			cfg.HTTP.Port = p
		}
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// validate は読み込んだ設定を検証する
func validate(c *Config) error {
	// JWT 秘密鍵の強度
	if len(c.Auth.JWTSecret) < 16 {
		return errors.New("JWT_SECRET too short (>=16)")
	}
	// アクセストークン有効期限
	if c.Auth.AccessExpiresIn <= 0 || c.Auth.AccessExpiresIn > 24*time.Hour {
		return errors.New("JWT_ACCESS_EXPIRES_IN out of range (0<ttl<=24h)")
	}
	// 本番 × '*'（AllowCredsとの整合もブラウザ仕様的にNG）
	for _, o := range c.HTTP.CORS.AllowOrigins {
		if o == "*" && c.HTTP.Env == "production" {
			return errors.New("CORS_ALLOW_ORIGINS must not be * in production")
		}
	}
	return nil
}
