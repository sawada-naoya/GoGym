// 役割: セッションドメインのエンティティ（認証・セッション管理）
// 受け取り: JTI, UserID, 有効期限, 現在時刻
// 処理: RefreshTokenの作成、バリデーション、有効性チェック
// 返却: 検証済みRefreshTokenエンティティ、エラー
package domain

import (
	"errors"
	"time"
)

// RefreshToken リフレッシュトークンエンティティ
type RefreshToken struct {
	JTI       string     // JWT ID (ULID)
	UserID    string     // ユーザーID (ULID)
	RevokedAt *time.Time // 無効化タイムスタンプ
	ExpiresAt time.Time  // 有効期限タイムスタンプ
	CreatedAt time.Time  // 作成タイムスタンプ
}

// NewRefreshToken バリデーション付きで新しいリフレッシュトークンを作成
func NewRefreshToken(jti, userID string, expiresAt, now time.Time) (*RefreshToken, error) {
	if jti == "" {
		return nil, errors.New("invalid jti")
	}
	if userID == "" {
		return nil, errors.New("invalid user id")
	}
	if !expiresAt.After(now) {
		return nil, errors.New("invalid expires at")
	}

	return &RefreshToken{
		JTI:       jti,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}, nil
}

// IsExpired 指定時刻でトークンが期限切れかチェック
func (rt *RefreshToken) IsExpired(at time.Time) bool {
	return !at.Before(rt.ExpiresAt)
}

// IsRevoked トークンが無効化されているかチェック
func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

// Revoke 指定時刻でトークンを無効化
func (rt *RefreshToken) Revoke(at time.Time) {
	rt.RevokedAt = &at
}

// IsValid 指定時刻でトークンが有効かチェック（期限切れでなく無効化もされていない）
func (rt *RefreshToken) IsValid(at time.Time) bool {
	return !rt.IsExpired(at) && !rt.IsRevoked()
}
