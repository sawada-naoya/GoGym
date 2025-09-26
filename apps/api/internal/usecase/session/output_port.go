// 役割: セッションユースケースの出力ポート（外部依存関係）
// 受け取り: 各種サービスの依存性注入インターフェース
// 処理: JWT生成、RefreshToken永続化、ID生成の抽象化
// 返却: 各サービスの実装への参照
package session

import (
	"context"
	"time"

	jwt "gogym-api/internal/adapter/auth"
	userDom "gogym-api/internal/domain/user"
)

type JWT interface {
	IssueAccess(userID string) (token string, ttl time.Duration, err error)
	IssueRefresh(userID string) (token string, ttl time.Duration, jti string, exp time.Time, err error)
	ParseRefresh(tokenStr string) (jwt.RefreshClaims, error)
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email userDom.Email) (*userDom.User, error)
}

type IDProvider interface {
	NewJTI() string
}

type TimeProvider interface {
	Now() time.Time
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) error
}
