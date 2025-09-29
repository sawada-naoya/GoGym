package auth

import (
	"errors"
	"fmt"

	uc "gogym-api/internal/usecase/user"

	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHasher は bcrypt を使用したパスワードハッシュ化を担当
type BcryptPasswordHasher struct {
	cost int
}

// NewBcryptPasswordHasher はデフォルトコストを使って新しいハッシャーを作成
func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{
		cost: bcrypt.DefaultCost,
	}
}

// インターフェース実装の確認
var _ uc.PasswordHasher = (*BcryptPasswordHasher)(nil)

// HashPassword は平文パスワードをハッシュ化
func (h *BcryptPasswordHasher) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword はパスワードとハッシュを照合
func (h *BcryptPasswordHasher) VerifyPassword(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return err
		}
		return fmt.Errorf("password verify failed: %w", err)
	}
	return nil
}
