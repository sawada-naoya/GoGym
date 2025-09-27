// di/sets.go
// 役割: 依存性注入のためのプロバイダー関数を定義
package di

import (
	"gogym-api/configs"
)

// ProvideDatabaseConfig は設定からデータベース設定を提供
func ProvideDatabaseConfig(cfg *configs.Config) configs.DatabaseConfig {
	return cfg.Database
}