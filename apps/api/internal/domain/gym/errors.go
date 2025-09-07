// internal/domain/gym/errors.go  
// 役割: ジムドメイン固有のエラー定義
// ジム関連のビジネスルール違反やバリデーションエラーの定義
package gym

// Gym domain specific errors
var (
	ErrGymNotFound      = &DomainError{Code: ErrNotFound, Key: "gym_not_found", Message: "gym not found"}
	ErrTagNotFound      = &DomainError{Code: ErrNotFound, Key: "tag_not_found", Message: "tag not found"}
	ErrInvalidGymName   = &DomainError{Code: ErrInvalidInput, Key: "invalid_gym_name", Message: "invalid gym name"}
	ErrInvalidTagName   = &DomainError{Code: ErrInvalidInput, Key: "invalid_tag_name", Message: "invalid tag name"}
	ErrFavoriteExists   = &DomainError{Code: ErrAlreadyExists, Key: "favorite_exists", Message: "gym already in favorites"}
	ErrFavoriteNotFound = &DomainError{Code: ErrNotFound, Key: "favorite_not_found", Message: "favorite not found"}
)