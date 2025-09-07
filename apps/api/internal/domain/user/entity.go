// internal/domain/user/entity.go
// 役割: ユーザードメインのEntity/VO（Domain Layer）
// ビジネスルールと不変条件を持つ純粋なドメインオブジェクト。GORM/JSONタグは一切なし
package user

import (
	"strings"
	"time"
)

// User はユーザーの集約ルートを表す
type User struct {
	BaseEntity
	Email        Email
	PasswordHash string
	DisplayName  string `validate:"required,max=100"`
}


// NewUser は検証付きで新しいユーザーを作成する
func NewUser(email Email, displayName string) (*User, error) {
	if !email.IsValid() {
		return nil, NewDomainError(ErrInvalidEmail, "invalid_email", "invalid email format")
	}

	user := &User{
		Email:       email,
		DisplayName: strings.TrimSpace(displayName),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate はユーザーデータを検証する
func (u *User) Validate() error {
	if !u.Email.IsValid() {
		return NewDomainError(ErrInvalidEmail, "invalid_email", "invalid email format")
	}

	if err := u.validateDisplayName(); err != nil {
		return err
	}

	return nil
}

// validateDisplayName は表示名を検証する
func (u *User) validateDisplayName() error {
	if u.DisplayName == "" {
		return NewDomainError(ErrInvalidInput, "invalid_display_name", "display name is required")
	}

	if len(u.DisplayName) > 100 {
		return NewDomainError(ErrInvalidInput, "invalid_display_name", "display name too long")
	}

	return nil
}

// SetPasswordHash はハッシュ化されたパスワードを設定する
func (u *User) SetPasswordHash(hash string) {
	u.PasswordHash = hash
}

// RefreshToken はリフレッシュトークンエンティティを表す
type RefreshToken struct {
	ID        ID
	UserID    ID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	User      *User
}


// NewRefreshToken は新しいリフレッシュトークンを作成する
func NewRefreshToken(userID ID, tokenHash string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
}

// IsExpired はトークンが期限切れかどうかをチェックする
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}