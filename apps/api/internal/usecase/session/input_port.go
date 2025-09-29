// 役割: セッションユースケースの入力ポート（外部からのインターフェース）
// 受け取り: HTTPハンドラーからのリクエスト
// 処理: セッション管理のビジネスロジック定義
// 返却: セッション作成・更新・削除の結果
package session

import (
	"context"
	"gogym-api/internal/adapter/dto"
)

type SessionUseCase interface {
	Login(ctx context.Context, req dto.LoginRequest) error
	CreateSession(ctx context.Context, email string) (dto.TokenResponse, error)
	// Logout(ctx context.Context, refreshToken string) error
}
