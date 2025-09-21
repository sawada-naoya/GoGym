// 役割: RefreshTokenエンティティ（Userとは分離）
// 時刻は引数で受け取り、ドメインからtime.Now()を排除
package user

import "time"

type RefreshToken struct {
	ID        ID
	UserID    ID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	User      *User // ユーザー情報（JOIN時に使用）
}

// NewRefreshToken: 不変条件を満たすトークンを生成
func NewRefreshToken(id ID, userID ID, tokenHash string, expiresAt time.Time, now time.Time) (*RefreshToken, error) {
	if tokenHash == "" {
		return nil, NewDomainError(ErrInvalidInput, "invalid_token_hash", "token hash required")
	}
	if !expiresAt.After(now) {
		return nil, NewDomainError(ErrInvalidInput, "invalid_expires_at", "expiresAt must be in the future")
	}
	return &RefreshToken{
		ID:        id,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}, nil
}

// IsExpired: 指定時刻時点で期限切れか判定
func (rt *RefreshToken) IsExpired(at time.Time) bool {
	return !at.Before(rt.ExpiresAt)
}
