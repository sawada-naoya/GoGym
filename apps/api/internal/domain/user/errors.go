// internal/domain/user/errors.go
// 役割: ユーザードメイン固有のエラー定義
// ユーザー認証・認可関連のビジネスルール違反やバリデーションエラーの定義
package user

import "gogym-api/internal/domain/common"

// User domain specific errors
var (
	ErrUserNotFound        = &common.DomainError{Code: common.ErrNotFound, Key: "user_not_found", Message: "user not found"}
	ErrUserAlreadyExists   = &common.DomainError{Code: common.ErrAlreadyExists, Key: "user_already_exists", Message: "user already exists"}
	ErrInvalidEmail        = &common.DomainError{Code: common.ErrInvalidEmail, Key: "invalid_email", Message: "invalid email format"}
	ErrWeakPassword        = &common.DomainError{Code: common.ErrWeakPassword, Key: "weak_password", Message: "password does not meet requirements"}
	ErrInvalidCredentials  = &common.DomainError{Code: common.ErrUnauthorized, Key: "invalid_credentials", Message: "invalid email or password"}
	ErrTokenExpired        = &common.DomainError{Code: common.ErrUnauthorized, Key: "token_expired", Message: "refresh token expired"}
	ErrInvalidToken        = &common.DomainError{Code: common.ErrUnauthorized, Key: "invalid_token", Message: "invalid token"}
)