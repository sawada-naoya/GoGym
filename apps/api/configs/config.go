package configs

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	Auth     AuthConfig
	Server   ServerConfig
	S3       S3Config
	Redis    RedisConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Timezone string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret        string
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
	Issuer           string
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port int
	Env  string
	CORS CORSConfig
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

// S3Config holds S3/MinIO configuration
type S3Config struct {
	Endpoint       string
	Bucket         string
	AccessKey      string
	SecretKey      string
	Region         string
	PublicURL      string
	ForcePathStyle bool
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// NewConfig creates application configuration from environment variables
func NewConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "gogym"),
			Password: getEnv("DB_PASSWORD", "password"),
			Database: getEnv("DB_NAME", "gogym"),
			Timezone: getEnv("TZ", "Asia/Tokyo"),
		},
		Auth: AuthConfig{
			JWTSecret:        getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			AccessExpiresIn:  parseDuration("JWT_ACCESS_EXPIRES_IN", "1h"),
			RefreshExpiresIn: parseDuration("JWT_REFRESH_EXPIRES_IN", "24h"),
			Issuer:           getEnv("JWT_ISSUER", "gogym-api"),
		},
		Server: ServerConfig{
			Port: parseInt("PORT", 8080),
			Env:  getEnv("APP_ENV", "development"),
			CORS: CORSConfig{
				AllowOrigins: []string{
					getEnv("NEXT_PUBLIC_WEB_URL", "http://localhost:3000"),
					"http://localhost:3000", // Development fallback
				},
				AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders: []string{"Content-Type", "Authorization"},
			},
		},
		S3: S3Config{
			Endpoint:       getEnv("S3_ENDPOINT", "http://localhost:9000"),
			Bucket:         getEnv("S3_BUCKET", "gogym-photos"),
			AccessKey:      getEnv("S3_ACCESS_KEY", "minioadmin"),
			SecretKey:      getEnv("S3_SECRET_KEY", "minioadmin123"),
			Region:         getEnv("S3_REGION", "us-east-1"),
			PublicURL:      getEnv("S3_PUBLIC_URL", "http://localhost:9000/gogym-photos"),
			ForcePathStyle: true, // MinIO requires path-style
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       parseInt("REDIS_DB", 0),
		},
	}
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func parseDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	// Parse default value
	parsed, _ := time.ParseDuration(defaultValue)
	return parsed
}
