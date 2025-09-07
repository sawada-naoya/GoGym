// internal/domain/gym/types.go  
// 役割: ジムドメインの基本型（Domain Layer）
// ジムドメイン固有の基本型、Location、検索関連型の定義。GORM/JSONタグは一切なし
package gym

import (
	"fmt"
	"time"
)

// ID はエンティティの一意識別子を表す
type ID int64

// BaseEntity はすべてのエンティティに共通するフィールドを提供する
type BaseEntity struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Location は地理座標を表す
type Location struct {
	Latitude  float64 `validate:"required,latitude"`
	Longitude float64 `validate:"required,longitude"`
}

// IsValid は位置座標を検証する
func (l Location) IsValid() bool {
	return l.Latitude >= -90 && l.Latitude <= 90 &&
		l.Longitude >= -180 && l.Longitude <= 180
}

// String は位置の文字列表現を返す
func (l Location) String() string {
	return fmt.Sprintf("(%f, %f)", l.Latitude, l.Longitude)
}

// SearchQuery は検索パラメータを表す
type SearchQuery struct {
	Query      string
	Location   *Location
	RadiusM    *int
	Pagination Pagination
}

// Pagination はページングパラメータを表す
type Pagination struct {
	Cursor string
	Limit  int
}

// PaginatedResult はページングされたレスポンスを表す
type PaginatedResult[T any] struct {
	Items      []T
	NextCursor string
	HasMore    bool
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
	ErrInvalidLocation
	ErrInvalidRating
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