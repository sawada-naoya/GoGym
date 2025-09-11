// internal/domain/user/errors.go
// 役割: ユーザードメイン固有のエラー定義
// ユーザー認証・認可関連のビジネスルール違反やバリデーションエラーの定義
package user

// User domain specific errors  
var (
	ErrUserNotFound        = &DomainError{Code: ErrNotFound, Key: "user_not_found", Message: "user not found"}
	ErrUserAlreadyExists   = &DomainError{Code: ErrAlreadyExists, Key: "user_already_exists", Message: "user already exists"}
	ErrInvalidEmailError   = &DomainError{Code: ErrInvalidEmail, Key: "invalid_email", Message: "invalid email format"}
	ErrWeakPasswordError   = &DomainError{Code: ErrWeakPassword, Key: "weak_password", Message: "password does not meet requirements"}
	ErrInvalidCredentials  = &DomainError{Code: ErrUnauthorized, Key: "invalid_credentials", Message: "invalid email or password"}
	ErrTokenExpired        = &DomainError{Code: ErrUnauthorized, Key: "token_expired", Message: "refresh token expired"}
	ErrInvalidToken        = &DomainError{Code: ErrUnauthorized, Key: "invalid_token", Message: "invalid token"}
)