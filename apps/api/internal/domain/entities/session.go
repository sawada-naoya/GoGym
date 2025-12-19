package domain

import (
	"errors"
	"time"
)

type RefreshToken struct {
	JTI       string     // JWT ID (ULID)
	UserID    string     // ユーザーID (ULID)
	RevokedAt *time.Time // 無効化タイムスタンプ
	ExpiresAt time.Time  // 有効期限タイムスタンプ
	CreatedAt time.Time  // 作成タイムスタンプ
}

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
