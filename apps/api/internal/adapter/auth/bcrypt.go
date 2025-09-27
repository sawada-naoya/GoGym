package auth

import (
	"errors"
	"fmt"

	uc "gogym-api/internal/usecase/user"

	"golang.org/x/crypto/bcrypt"
)

/*
役割: パスワードハッシュ化サービス（Auth Layer）
受け取り: ハッシュ化コスト、ペッパー文字列
処理: bcryptを使用したパスワードのハッシュ化と検証
返却: ハッシュ化されたパスワード、検証結果
使用方法:
  // ハッシャー初期化
  hasher, _ := auth.NewBcryptPasswordHasher(12, "secret-pepper")

  // パスワードハッシュ化（ユーザー登録時）
  hashedPassword, _ := hasher.HashPassword("user-password")

  // パスワード検証（ログイン時）
  err := hasher.VerifyPassword("user-password", hashedPassword)
  if err == nil {
      // 認証成功
  }
*/

// BcryptPasswordHasher はbcryptを使用したパスワードハッシュ化を担当
type BcryptPasswordHasher struct {
	cost   int    // ハッシュ化の計算コスト（デフォルト推奨）
	pepper []byte // セキュリティ強化のための追加文字列（オプション）
}

// NewBcryptPasswordHasher は新しいbcryptハッシャーを作成
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

// インターフェース実装の確認
var _ uc.PasswordHasher = (*BcryptPasswordHasher)(nil)

// HashPassword は平文パスワードをハッシュ化
func (h *BcryptPasswordHasher) HashPassword(password string) (string, error) {
	// パスワード + pepper を結合（例: "abc" + "d" -> []byte[97, 98, 99, 100]）
	pw := append([]byte(password), h.pepper...)
	hash, err := bcrypt.GenerateFromPassword(pw, h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword はパスワードとハッシュを照合
func (h *BcryptPasswordHasher) VerifyPassword(password, hash string) error {
	// パスワード + pepper を結合（例: "abc" + "d" -> []byte[97, 98, 99, 100]）
	pw := append([]byte(password), h.pepper...)
	// ハッシュとパスワードを照合
	if err := bcrypt.CompareHashAndPassword([]byte(hash), pw); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return err
		}
		return fmt.Errorf("password verify failed: %w", err)
	}
	return nil
}
