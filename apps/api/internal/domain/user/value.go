// Domain Value Objects: transport/GORM/ctx依存なし
package user

import (
	"regexp"
	"strings"
)

var (
	emailRe          = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
	maxEmailLen      = 320 //メールアドレスの最大長
	maxLocalLen      = 64  // @より前の部分の最大長
	maxDomainPartLen = 255 // @の後のドメイン部分の最大長
)

type Email struct {
	value string
}

// NewEmail: 小文字化 + trim + 形式検証
func NewEmail(raw string) (Email, error) {
	n := strings.ToLower(strings.TrimSpace(raw))
	if n == "" {
		return Email{}, NewDomainError("invalid_email")
	}

	if len(n) > maxEmailLen {
		return Email{}, NewDomainError("invalid_email")
	}

	if !emailRe.MatchString(n) {
		return Email{}, NewDomainError("invalid_email")
	}

	// @がない、または@の前後が空文字列の場合は無効
	// 例: "@example.com", "user@"
	if at := strings.LastIndexByte(n, '@'); at <= 0 || at == len(n)-1 {
		return Email{}, NewDomainError("invalid_email")
	} else {
		local := n[:at]
		domain := n[at+1:]
		if len(local) > maxLocalLen || len(domain) > maxDomainPartLen {
			return Email{}, NewDomainError("invalid_email")
		}
	}

	return Email{value: n}, nil
}

// String はEmailの文字列表現を返す
func (e Email) String() string {
	return e.value
}
