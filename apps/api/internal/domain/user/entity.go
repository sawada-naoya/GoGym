// 役割: ユーザードメインのEntity/VO（Domain Layer）
// ビジネスルールと不変条件を持つ純粋なドメインオブジェクト。GORM/JSONタグは一切なし
package user

import (
	"strings"
	"time"
)

type User struct {
	ID           string    // ULID識別子
	Name         string    // 表示名
	Email        Email     // メールアドレス（バリューオブジェクト）
	PasswordHash string    // パスワードハッシュ
	CreatedAt    time.Time // 作成日時
	UpdatedAt    time.Time // 更新日時
}

// NewUser: 不変条件を満たすユーザーを生成（IDは自動生成）
func NewUser(id string, name string, email Email, passwordHash string, now time.Time) (*User, error) {
	n := strings.TrimSpace(name)
	if n == "" || len(n) > 100 {
		return nil, NewDomainError(ErrInvalidInput, "invalid_name", "name required and <=100 chars")
	}
	if strings.TrimSpace(passwordHash) == "" {
		return nil, NewDomainError(ErrInvalidInput, "invalid_password_hash", "password hash required")
	}

	return &User{
		ID:           id,
		Name:         n,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Rename: 表示名を変更（不変条件を再度満たすこと）
func (u *User) Rename(newName string) error {
	n := strings.TrimSpace(newName)
	if n == "" || len(n) > 100 {
		return NewDomainError(ErrInvalidInput, "invalid_name", "name required and <=100 chars")
	}
	u.Name = n
	u.UpdatedAt = time.Now() // 更新時刻を更新
	return nil
}

// RotatePasswordHash: パスワードハッシュを安全に更新（rawは扱わない）
func (u *User) RotatePasswordHash(newHash string) error {
	if strings.TrimSpace(newHash) == "" {
		return NewDomainError(ErrInvalidInput, "invalid_password_hash", "password hash required")
	}
	u.PasswordHash = newHash
	u.UpdatedAt = time.Now() // 更新時刻を更新
	return nil
}
