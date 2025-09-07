// internal/domain/common/errors.go
// 役割: 共通ドメインエラー定義
// 全ドメインで共有されるエラー型とエラーコードの定義
package common

import "fmt"

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
	ErrInvalidLocation
	ErrInvalidRating
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