package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

type User struct {
	ID           ulid.ULID // ULID識別子
	Name         string    // 表示名
	Email        string    // メールアドレス（バリューオブジェクト）
	PasswordHash string    // パスワードハッシュ
	CreatedAt    time.Time // 作成日時
	UpdatedAt    time.Time // 更新日時
}

func NewUser(id ulid.ULID, name, email, passwordHash string, now time.Time) *User {
	n := strings.TrimSpace(name)
	if n == "" || len(n) > 100 {
		return nil
	}
	if strings.TrimSpace(passwordHash) == "" {
		return nil
	}

	return &User{
		ID:           id,
		Name:         n,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Rename: 表示名を変更（不変条件を再度満たすこと）
func (u *User) Rename(newName string) error {
	n := strings.TrimSpace(newName)
	if n == "" || len(n) > 100 {
		return errors.New("invalid name")
	}
	u.Name = n
	u.UpdatedAt = time.Now() // 更新時刻を更新
	return nil
}

// RotatePasswordHash: パスワードハッシュを安全に更新（rawは扱わない）
func (u *User) RotatePasswordHash(newHash string) error {
	if strings.TrimSpace(newHash) == "" {
		return errors.New("invalid password hash")
	}
	u.PasswordHash = newHash
	u.UpdatedAt = time.Now() // 更新時刻を更新
	return nil
}
