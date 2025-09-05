package common

import "errors"

// Domain errors
var (
	ErrNotFound         = errors.New("resource not found")
	ErrAlreadyExists    = errors.New("resource already exists")
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrInternal         = errors.New("internal error")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrWeakPassword     = errors.New("password too weak")
	ErrInvalidLocation  = errors.New("invalid location coordinates")
	ErrInvalidRating    = errors.New("rating must be between 1 and 5")
)

// DomainError wraps errors with additional context
type DomainError struct {
	Err     error
	Code    string
	Message string
	Details map[string]interface{}
}

func (e DomainError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error
func NewDomainError(err error, code string, message string) DomainError {
	return DomainError{
		Err:     err,
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adds details to domain error
func (e DomainError) WithDetails(details map[string]interface{}) DomainError {
	e.Details = details
	return e
}