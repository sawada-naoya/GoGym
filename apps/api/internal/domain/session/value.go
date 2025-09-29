// 役割: セッションドメインのValue Object（値オブジェクト）
// 受け取り: アクセストークン、リフレッシュトークン、有効期限、セッション情報
// 処理: 不変のトークンペア作成、セッション情報の管理
// 返却: 不変なトークンペア、セッション情報オブジェクト
package session

import (
	"time"
)

// TokenPair アクセストークンとリフレッシュトークンのペア
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // 秒単位
}

// NewTokenPair 新しいトークンペアを作成
func NewTokenPair(accessToken, refreshToken string, expiresIn int) TokenPair {
	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
}

// SessionInfo 検証用のセッション情報
type SessionInfo struct {
	UserID    string
	JTI       string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// NewSessionInfo 新しいセッション情報を作成
func NewSessionInfo(userID, jti string, issuedAt, expiresAt time.Time) SessionInfo {
	return SessionInfo{
		UserID:    userID,
		JTI:       jti,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}
}

// IsExpired 指定時刻でセッションが期限切れかチェック
func (si SessionInfo) IsExpired(at time.Time) bool {
	return at.After(si.ExpiresAt)
}