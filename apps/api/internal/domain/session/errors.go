// 役割: セッションドメインのエラー定義
// 受け取り: エラーコード、キー、メッセージ、原因エラー
// 処理: ドメイン固有のエラー作成、エラー情報の整理
// 返却: 構造化されたドメインエラー
package session

import "fmt"

// DomainError セッションドメイン固有のエラー
type DomainError struct {
	Code    ErrorCode
	Key     string
	Message string
	Cause   error
}

// ErrorCode エラー分類
type ErrorCode int

const (
	ErrUnknown ErrorCode = iota
	ErrInvalidInput
	ErrNotFound
	ErrAlreadyExists
	ErrUnauthorized
	ErrTokenExpired
	ErrTokenRevoked
)

// NewDomainError エラーコード付きで新しいドメインエラーを作成
func NewDomainError(code ErrorCode, key, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Key:     key,
		Message: message,
		Cause:   nil,
	}
}

// NewDomainErrorWithCause 原因エラー付きで新しいドメインエラーを作成
func NewDomainErrorWithCause(cause error, key, message string) *DomainError {
	code := ErrUnknown
	if domainErr, ok := cause.(*DomainError); ok {
		code = domainErr.Code
	} else if cause != nil {
		code = ErrUnknown
	}

	return &DomainError{
		Code:    code,
		Key:     key,
		Message: message,
		Cause:   cause,
	}
}

// Error errorインターフェースを実装
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %s)", e.Key, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

// Unwrap 根本原因エラーを返却
func (e *DomainError) Unwrap() error {
	return e.Cause
}