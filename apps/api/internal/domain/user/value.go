// internal/domain/user/value.go
// 役割: ユーザードメインのValue Object（Domain Layer）
// 不変性と検証ロジックを持つ純粋なドメインバリューオブジェクト。GORM/JSONタグは一切なし
package user

import (
	"gogym-api/internal/domain/common"
	"regexp"
	"strings"
	"unicode"
)

// Email は検証と正規化機能付きのメールアドレスバリューオブジェクトを表す
type Email string

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail は検証と正規化付きで新しいメールアドレスを作成する
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

// IsValid はメールアドレスの形式を検証する
func (e Email) IsValid() bool {
	if e == "" || len(string(e)) > 255 {
		return false
	}
	return emailRegex.MatchString(string(e))
}

// String は文字列表現を返す
func (e Email) String() string {
	return string(e)
}

// Domain はメールアドレスのドメイン部分を返す
func (e Email) Domain() string {
	parts := strings.Split(string(e), "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// Password はパスワード検証ルールを表す
type Password struct {
	value string
}

// NewPassword は検証付きで新しいパスワードを作成する
func NewPassword(password string) (*Password, error) {
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}
	
	return &Password{value: password}, nil
}

// String はパスワード値を返す（慎重に使用）
func (p *Password) String() string {
	return p.value
}

// ValidatePassword はパスワード強度を検証する
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

// ValueError はバリューオブジェクト検証エラーを表す
type ValueError struct {
	Field   string
	Message string
}

func (e *ValueError) Error() string {
	return e.Field + ": " + e.Message
}