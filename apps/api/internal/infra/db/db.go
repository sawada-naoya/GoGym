// internal/infra/db/gorm.go
// 役割: データベース接続インフラストラクチャ（Infrastructure Layer）
// GORM（Go ORM）を使用してMySQL 8.0への接続とコネクションプールの設定を担当
package db

import (
	"fmt"
	"time"

	"gogym-api/internal/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB は設定情報を元にGORMのDBインスタンスを生成する
// 接続プールやログ設定も同時に行い、本番環境での運用に適した設定を適用する
//
// Parameters:
//   - cfg: データベース接続設定（ホスト、ポート、ユーザー、パスワード等）
//
// Returns:
//   - *gorm.DB: 設定済みのGORMデータベースインスタンス
//   - error: 接続エラーまたは設定エラー
func NewDB(cfg configs.DatabaseConfig) (*gorm.DB, error) {
	// MySQL接続用のDSN（Data Source Name）を構築
	// parseTime=true: time.Time型への自動変換を有効化
	// loc: タイムゾーン設定（URLエンコードが必要）
	// charset/collation: UTF-8（日本語対応）+ 大文字小文字区別なし
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci&interpolateParams=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	// GORM用ロガー設定
	// 開発環境では実行されるSQLクエリを全てログ出力する（デバッグ用）
	// 本番環境ではエラーのみ出力するよう後で調整可能
	newLogger := logger.Default.LogMode(logger.Info)

	// GORM DBインスタンスを作成
	// MySQL 8.0ドライバーを使用してデータベースに接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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
