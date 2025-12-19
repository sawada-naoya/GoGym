package db

import (
	"fmt"
	"os"
	"time"

	"gogym-api/internal/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(cfg configs.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.Port,
	)

	// 環境に応じたログレベルの設定
	var newLogger logger.Interface
	env := os.Getenv("APP_ENV")
	if env == "local" || env == "development" || env == "" {
		// ローカル環境ではInfoレベルで詳細なログを表示
		newLogger = logger.Default.LogMode(logger.Info)
	} else {
		// 本番環境ではWarnレベル以上のみ表示
		newLogger = logger.Default.LogMode(logger.Warn)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("データベース接続に失敗しました: %w", err)
	}

	// 接続プール設定のため、内部のsql.DBインスタンスを取得
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("接続プール設定の取得に失敗しました: %w", err)
	}

	// 接続プールの詳細設定
	// 本番環境でのパフォーマンスと安定性を考慮した設定値
	sqlDB.SetMaxOpenConns(25)                 // 最大同時接続数（高負荷時でもリソース枯渇を防ぐ）
	sqlDB.SetMaxIdleConns(10)                 // アイドル状態で保持する接続数（レスポンス速度向上）
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 接続の最大生存時間（接続リークを防ぐ）
	sqlDB.SetConnMaxIdleTime(1 * time.Minute) // アイドル接続の最大生存時間

	return db, nil
}
