package di

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"

	"gogym-api/internal/adapter/auth"
	userUC "gogym-api/internal/usecase/user"
)

// NewPasswordHasher はBcryptPasswordHasherのインスタンスを作成
func NewPasswordHasher() userUC.PasswordHasher {
	hasher, _ := auth.NewBcryptPasswordHasher(0, "") // デフォルト設定
	return hasher
}

// ULIDProvider はULIDを生成するIDProviderの実装
type ULIDProvider struct{}

// NewIDProvider はULIDProviderのインスタンスを作成
func NewIDProvider() userUC.IDProvider {
	return &ULIDProvider{}
}

// NewUserID は新しいULIDを生成
func (p *ULIDProvider) NewUserID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
}