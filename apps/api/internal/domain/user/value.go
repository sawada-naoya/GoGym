// internal/domain/user/value.go
// 役割: ユーザードメインのバリューオブジェクト
// Email正規化/検証、Password強度チェック等のユーザー関連バリューオブジェクトの定義
package user

import (
	"gogym-api/internal/domain/common"
	"regexp"
	"strings"
	"unicode"
)

// Email represents email value object with validation and normalization
type Email string

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new email with validation and normalization
func NewEmail(email string) (Email, error) {
	normalized := strings.TrimSpace(strings.ToLower(email))
	
	if normalized == "" {
		return "", &ValueError{Field: "email", Message: "email is required"}
	}

	if len(normalized) > 255 {
		return "", &ValueError{Field: "email", Message: "email too long"}
	}

	if !emailRegex.MatchString(normalized) {
		return "", &ValueError{Field: "email", Message: "invalid email format"}
	}

	return Email(normalized), nil
}

// IsValid validates the email format
func (e Email) IsValid() bool {
	if e == "" || len(string(e)) > 255 {
		return false
	}
	return emailRegex.MatchString(string(e))
}

// String returns string representation
func (e Email) String() string {
	return string(e)
}

// Domain returns the domain part of the email
func (e Email) Domain() string {
	parts := strings.Split(string(e), "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// Password represents password validation rules
type Password struct {
	value string
}

// NewPassword creates a new password with validation
func NewPassword(password string) (*Password, error) {
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}
	
	return &Password{value: password}, nil
}

// String returns the password value (use carefully)
func (p *Password) String() string {
	return p.value
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

// ValueError represents value object validation error
type ValueError struct {
	Field   string
	Message string
}

func (e *ValueError) Error() string {
	return e.Field + ": " + e.Message
}