package user

import (
	"gogym-api/internal/domain/common"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// User represents a user entity
type User struct {
	common.BaseEntity
	Email        string    `json:"email" gorm:"uniqueIndex;size:255;not null" validate:"required,email"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	DisplayName  string    `json:"display_name" gorm:"size:100;not null" validate:"required,max=100"`
}

// TableName returns the table name for GORM
func (User) TableName() string {
	return "users"
}

// NewUser creates a new user with validation
func NewUser(email, displayName string) (*User, error) {
	user := &User{
		Email:       strings.TrimSpace(strings.ToLower(email)),
		DisplayName: strings.TrimSpace(displayName),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate validates user data
func (u *User) Validate() error {
	if err := u.validateEmail(); err != nil {
		return err
	}

	if err := u.validateDisplayName(); err != nil {
		return err
	}

	return nil
}

// validateEmail validates email format
func (u *User) validateEmail() error {
	if u.Email == "" {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_email", "email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return common.NewDomainError(common.ErrInvalidEmail, "invalid_email", "invalid email format")
	}

	if len(u.Email) > 255 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_email", "email too long")
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

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return common.NewDomainError(common.ErrWeakPassword, "weak_password", "password must be at least 8 characters")
	}

	if len(password) > 128 {
		return common.NewDomainError(common.ErrWeakPassword, "weak_password", "password too long")
	}

	var hasUpper, hasLower, hasNumber bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsNumber(r):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return common.NewDomainError(
			common.ErrWeakPassword,
			"weak_password",
			"password must contain uppercase, lowercase, and number",
		)
	}

	return nil
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

// IsExpired checks if the token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}