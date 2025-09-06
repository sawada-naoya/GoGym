// apps/api/internal/adapter/db/gorm/provider.go
// GORM関連の依存性注入（DI）プロバイダーを定義するファイル
// Google Wireを使用してデータベース接続とリポジトリの生成を自動化する
package gorm

import (
	"gogym-api/configs"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// NewGormDB はWire DIコンテナ用のGORMデータベースプロバイダー
// configs.DatabaseConfigから*gorm.DBインスタンスを生成する
//
// この関数は実際のアプリケーション起動時に一度だけ実行され、
// 生成されたDBインスタンスは全てのリポジトリで共有される
func NewGormDB(cfg configs.DatabaseConfig) (*gorm.DB, error) {
	return NewDB(cfg)
}

// GormProviderSet はGORM関連の全てのプロバイダーをまとめたWire Set
// 他のパッケージからインポートしてDIコンテナに組み込む
var GormProviderSet = wire.NewSet(
	NewGormDB, // データベース接続プロバイダー
	// 今後、各リポジトリのプロバイダーもここに追加していく
	// NewUserRepository,
	// NewGymRepository,
	// NewTagRepository,
)