// 役割: 時刻取得サービス（Service Layer）
// 受け取り: なし
// 処理: 現在時刻の取得（テスト時にモック可能）
// 返却: 現在時刻
package service

import (
	"time"

	sessionUC "gogym-api/internal/usecase/session"
)

// timeProvider は session.TimeProvider インターフェースを実装
type timeProvider struct{}

// NewTimeProvider は新しいTimeProviderを作成
func NewTimeProvider() sessionUC.TimeProvider {
	return &timeProvider{}
}

// Now は現在時刻を返却
func (tp *timeProvider) Now() time.Time {
	return time.Now()
}