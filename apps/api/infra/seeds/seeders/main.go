// GoGym データベースシーダー
// NDJSON形式のシードデータを読み込んでデータベースに投入する
// 使用方法: go run infra/seeds/seeders/main.go
package main

import (
	"context"
	"log/slog"
	"os"

	"gogym-api/configs"
	gormdb "gogym-api/internal/adapter/db/gorm"
)

func main() {
	// ロガー初期化
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	// 設定読み込み
	cfg, err := configs.Load()
	if err != nil {
		logger.Error("設定の読み込みに失敗しました", "error", err)
		os.Exit(1)
	}

	// データベース接続
	db, err := gormdb.NewDB(cfg.Database)
	if err != nil {
		logger.Error("データベース接続に失敗しました", "error", err)
		os.Exit(1)
	}

	logger.Info("シーダー実行開始")

	ctx := context.Background()
	
	// シーダー実行順序（外部キー制約を考慮）
	seeders := []Seeder{
		NewUserSeeder(db, logger),
		NewGymSeeder(db, logger),
		NewReviewSeeder(db, logger),
	}

	for _, seeder := range seeders {
		if err := seeder.Run(ctx); err != nil {
			logger.Error("シーダー実行エラー", "seeder", seeder.Name(), "error", err)
			os.Exit(1)
		}
		logger.Info("シーダー完了", "seeder", seeder.Name())
	}

	logger.Info("全シーダー実行完了")
}

// Seeder インターフェース
type Seeder interface {
	Name() string
	Run(ctx context.Context) error
}