package log

import (
	"log/slog"
	"os"
)

func NewLogger(env string) *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: env == "development",
	}

	var h slog.Handler
	if env == "development" {
		opts.Level = slog.LevelDebug // 開発環境ではDEBUGレベルも表示
		h = slog.NewTextHandler(os.Stdout, opts)
	} else {
		opts.Level = slog.LevelInfo // 本番環境はINFO以上
		h = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(h).With("service", "gogym-api")
}

// IsDev は開発環境かどうかを判定
func IsDev() bool {
	return os.Getenv("APP_ENV") == "development"
}
