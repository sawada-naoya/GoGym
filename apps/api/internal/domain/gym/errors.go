// internal/domain/gym/errors.go  
// 役割: ジムドメイン固有のエラー定義
// ジム関連のビジネスルール違反やバリデーションエラーの定義
package gym

import "gogym-api/internal/domain/common"

// Gym domain specific errors
var (
	ErrGymNotFound      = &common.DomainError{Code: common.ErrNotFound, Key: "gym_not_found", Message: "gym not found"}
	ErrTagNotFound      = &common.DomainError{Code: common.ErrNotFound, Key: "tag_not_found", Message: "tag not found"}
	ErrInvalidGymName   = &common.DomainError{Code: common.ErrInvalidInput, Key: "invalid_gym_name", Message: "invalid gym name"}
	ErrInvalidTagName   = &common.DomainError{Code: common.ErrInvalidInput, Key: "invalid_tag_name", Message: "invalid tag name"}
	ErrFavoriteExists   = &common.DomainError{Code: common.ErrAlreadyExists, Key: "favorite_exists", Message: "gym already in favorites"}
	ErrFavoriteNotFound = &common.DomainError{Code: common.ErrNotFound, Key: "favorite_not_found", Message: "favorite not found"}
)