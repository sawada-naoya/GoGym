// internal/domain/user/model.go
// 役割: ユーザードメインのエンティティモデル
// ユーザー集約ルート: UserID, Email(VO), DisplayName の管理
package user

import (
	"gogym-api/internal/domain/common"
	"strings"
	"time"
)

// User represents the user aggregate root
type User struct {
	common.BaseEntity
	Email        Email     `json:"email" gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	DisplayName  string    `json:"display_name" gorm:"size:100;not null" validate:"required,max=100"`
}

// TableName returns the table name for GORM
func (User) TableName() string {
	return "users"
}

// NewUser creates a new user with validation
func NewUser(email Email, displayName string) (*User, error) {
	if !email.IsValid() {
		return nil, common.NewDomainError(common.ErrInvalidEmail, "invalid_email", "invalid email format")
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

// Validate validates user data
func (u *User) Validate() error {
	if !u.Email.IsValid() {
		return common.NewDomainError(common.ErrInvalidEmail, "invalid_email", "invalid email format")
	}

	if err := u.validateDisplayName(); err != nil {
		return err
	}

	return nil
}

// validateDisplayName validates display name
func (u *User) validateDisplayName() error {
	if u.DisplayName == "" {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_display_name", "display name is required")
	}

	if len(u.DisplayName) > 100 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_display_name", "display name too long")
	}

	return nil
}

// SetPasswordHash sets the hashed password
func (u *User) SetPasswordHash(hash string) {
	u.PasswordHash = hash
}

// RefreshToken represents a refresh token entity
type RefreshToken struct {
	ID        common.ID `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    common.ID `json:"user_id" gorm:"not null;index"`
	TokenHash string    `json:"-" gorm:"uniqueIndex;size:255;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for GORM
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// NewRefreshToken creates a new refresh token
func NewRefreshToken(userID common.ID, tokenHash string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
}

// IsExpired checks if the token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}