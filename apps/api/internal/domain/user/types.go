// internal/domain/user/types.go
// 役割: ユーザードメインの基本型（Domain Layer）
// ユーザードメイン固有の基本型とエラー定義。GORM/JSONタグは一切なし
package user

import (
	"crypto/rand"
	"fmt"
	"time"
	"github.com/oklog/ulid/v2"
)

// ID はULIDベースのエンティティ識別子
type ID string

// GenerateID は新しいULIDを生成する
func GenerateID() ID {
	return ID(ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String())
}

// ParseID はULID文字列からIDを作成する（バリデーション付き）
func ParseID(s string) (ID, error) {
	if _, err := ulid.Parse(s); err != nil {
		return "", NewDomainError(ErrInvalidInput, "invalid_id", "invalid ULID format")
	}
	return ID(s), nil
}

// DomainError represents a domain-specific error
type DomainError struct {
	Code    ErrorCode
	Key     string
	Message string
	Cause   error
}

// ErrorCode represents error classification
type ErrorCode int

const (
	ErrUnknown ErrorCode = iota
	ErrInvalidInput
	ErrNotFound
	ErrAlreadyExists
	ErrUnauthorized
	ErrForbidden
	ErrInternal
	ErrInvalidEmail
	ErrWeakPassword
)

// NewDomainError creates a new domain error with error code
func NewDomainError(code ErrorCode, key, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Key:     key,
		Message: message,
		Cause:   nil,
	}
}

// NewDomainErrorWithCause creates a new domain error with cause
func NewDomainErrorWithCause(cause error, key, message string) *DomainError {
	code := ErrUnknown
	if domainErr, ok := cause.(*DomainError); ok {
		code = domainErr.Code
	} else if cause != nil {
		code = ErrInternal
	}

	return &DomainError{
		Code:    code,
		Key:     key,
		Message: message,
		Cause:   cause,
	}
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %s)", e.Key, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

// Unwrap returns the underlying cause
func (e *DomainError) Unwrap() error {
	return e.Cause
}