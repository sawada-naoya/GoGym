// Domain Value Objects: transport/GORM/ctx依存なし
package user

import (
	"regexp"
	"strings"
)

// ===== Email (安全なVO) =====

// NOTE: RFCは長いので現実的フォーマット + 正規化のみ
var (
	emailRe          = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
	maxEmailLen      = 320 // 実務上の上限
	maxLocalLen      = 64  // ローカル部の目安
	maxDomainPartLen = 255 // ドメイン部の目安
)

// Emailは不変条件を満たす値のみ生成できるよう非公開フィールドにする
type Email struct {
	value string
}

// NewEmail: 小文字化 + trim + 形式検証
func NewEmail(raw string) (Email, error) {
	n := strings.ToLower(strings.TrimSpace(raw))
	if n == "" {
		return Email{}, NewDomainError(ErrInvalidEmail, "invalid_email", "email is required")
	}
	if len(n) > maxEmailLen {
		return Email{}, NewDomainError(ErrInvalidEmail, "invalid_email", "email too long")
	}
	if !emailRe.MatchString(n) {
		return Email{}, NewDomainError(ErrInvalidEmail, "invalid_email", "invalid email format")
	}
	// 参考程度の長さ検査（厳密RFC準拠はやりすぎ）
	if at := strings.LastIndexByte(n, '@'); at <= 0 || at == len(n)-1 {
		return Email{}, NewDomainError(ErrInvalidEmail, "invalid_email", "invalid email format")
	} else {
		local := n[:at]
		domain := n[at+1:]
		if len(local) > maxLocalLen || len(domain) > maxDomainPartLen {
			return Email{}, NewDomainError(ErrInvalidEmail, "invalid_email", "email part too long")
		}
	}

	return Email{value: n}, nil
}

// String: 正規化済み文字列を返す
func (e Email) String() string { return e.value }

// Domain: ドメイン部を返す（なければ空文字）
func (e Email) Domain() string {
	if i := strings.LastIndexByte(e.value, '@'); i >= 0 && i+1 < len(e.value) {
		return e.value[i+1:]
	}
	return ""
}

// Equals: 大文字小文字非区別の比較（正規化済みなので通常の==でも同値）
func (e Email) Equals(other Email) bool { return e.value == other.value }
