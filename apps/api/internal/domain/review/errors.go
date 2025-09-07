// internal/domain/review/errors.go
// 役割: レビュードメイン固有のエラー定義
// レビュー関連のビジネスルール違反やバリデーションエラーの定義
package review

import "gogym-api/internal/domain/common"

// Review domain specific errors
var (
	ErrReviewNotFound      = &common.DomainError{Code: common.ErrNotFound, Key: "review_not_found", Message: "review not found"}
	ErrInvalidRating       = &common.DomainError{Code: common.ErrInvalidRating, Key: "invalid_rating", Message: "rating must be between 1 and 5"}
	ErrInvalidPhotoURL     = &common.DomainError{Code: common.ErrInvalidInput, Key: "invalid_photo_url", Message: "invalid photo URL"}
	ErrReviewAlreadyExists = &common.DomainError{Code: common.ErrAlreadyExists, Key: "review_already_exists", Message: "user already reviewed this gym"}
	ErrUnauthorizedReview  = &common.DomainError{Code: common.ErrUnauthorized, Key: "unauthorized_review", Message: "user not authorized to modify this review"}
)