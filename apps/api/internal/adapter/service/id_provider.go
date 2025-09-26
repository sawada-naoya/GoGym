// 役割: ID生成サービス（Service Layer）
// 受け取り: なし（各種ULIDを生成）
// 処理: ユーザーID、JWT ID（JTI）等の一意識別子生成
// 返却: ULID文字列
package service

import (
	crypto "crypto/rand"
	"time"

	sessionUC "gogym-api/internal/usecase/session"
	userUC "gogym-api/internal/usecase/user"
	"github.com/oklog/ulid/v2"
)

/*
使用方法:
  provider := service.NewIDProvider()
  userID := provider.NewUserID()    // ユーザー登録時
  jti := provider.NewJTI()          // JWT生成時
*/

// IDProvider は各種ID生成を担当
type IDProvider struct {
	entropy *ulid.MonotonicEntropy
}

// NewIDProvider は新しいIDProviderを作成
func NewIDProvider() *IDProvider {
	return &IDProvider{
		entropy: ulid.Monotonic(crypto.Reader, 0),
	}
}

// NewUserID は新しいユーザーID（ULID）を生成
func (p *IDProvider) NewUserID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), p.entropy).String()
}

// NewJTI は新しいJWT ID（ULID）を生成
func (p *IDProvider) NewJTI() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), p.entropy).String()
}

// インターフェース実装の確認
var _ userUC.IDProvider = (*IDProvider)(nil)
var _ sessionUC.IDProvider = (*IDProvider)(nil)