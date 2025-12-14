// 役割: セッションユースケースの出力ポート（外部依存関係）
// 受け取り: 各種サービスの依存性注入インターフェース
// 処理: JWT生成、RefreshToken永続化、ID生成の抽象化
// 返却: 各サービスの実装への参照
package session

import (
	"context"

	dom "gogym-api/internal/domain/entities"
	"github.com/oklog/ulid/v2"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email dom.Email) (*dom.User, error)
	FindByID(ctx context.Context, id ulid.ULID) (*dom.User, error)
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) error
}
