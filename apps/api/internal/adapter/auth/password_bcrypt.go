package auth

import (
	"errors"
	"fmt"

	uc "gogym-api/internal/usecase/user"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct {
	cost   int    // costは基本デフォルトを使用する
	pepper []byte // pepperはセキュリティ強化のための追加文字列(なくてもいいよ)
}

func NewBcryptPasswordHasher(cost int, pepper string) (*BcryptPasswordHasher, error) {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	if cost < 4 || cost > 31 {
		return nil, fmt.Errorf("bcrypt cost must be between 4 and 31, got %d", cost)
	}

	return &BcryptPasswordHasher{
		cost:   cost,
		pepper: []byte(pepper),
	}, nil
}

// BcryptPasswordHasher構造体がuc.PasswordHasherインターフェースを実装していることを確認
var _ uc.PasswordHasher = (*BcryptPasswordHasher)(nil)

func (h *BcryptPasswordHasher) HashPassword(password string) (string, error) {
	//　例: abc + d -> []byte[97, 98, 99, 100]
	pw := append([]byte(password), h.pepper...)
	hash, err := bcrypt.GenerateFromPassword(pw, h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *BcryptPasswordHasher) VerifyPassword(password, hash string) error {
	pw := append([]byte(password), h.pepper...)
	if err := bcrypt.CompareHashAndPassword([]byte(hash), pw); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return err
		}
		return fmt.Errorf("password verify failed: %w", err)
	}
	return nil
}
