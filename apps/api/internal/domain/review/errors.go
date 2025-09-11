// internal/domain/review/errors.go
// 役割: レビュードメイン固有のエラー定義
// レビュー関連のビジネスルール違反やバリデーションエラーの定義
package review

import "gogym-api/internal/domain/gym"

// Review domain specific errors
var (
	ErrReviewNotFound      = &gym.DomainError{Code: gym.ErrNotFound, Key: "review_not_found", Message: "review not found"}
	ErrInvalidRating       = &gym.DomainError{Code: gym.ErrInvalidRating, Key: "invalid_rating", Message: "rating must be between 1 and 5"}
	ErrInvalidPhotoURL     = &gym.DomainError{Code: gym.ErrInvalidInput, Key: "invalid_photo_url", Message: "invalid photo URL"}
	ErrReviewAlreadyExists = &gym.DomainError{Code: gym.ErrAlreadyExists, Key: "review_already_exists", Message: "user already reviewed this gym"}
	ErrUnauthorizedReview  = &gym.DomainError{Code: gym.ErrUnauthorized, Key: "unauthorized_review", Message: "user not authorized to modify this review"}
)