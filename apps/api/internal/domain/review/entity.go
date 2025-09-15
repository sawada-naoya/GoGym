// internal/domain/review/entity.go
// 役割: レビュードメインのEntity/VO（Domain Layer）
// ビジネスルールと不変条件を持つ純粋なドメインオブジェクト。GORM/JSONタグは一切なし
package review

import (
	"gogym-api/internal/domain/gym"
	"strings"
)

// Review represents a gym review aggregate root
type Review struct {
	gym.BaseEntity
	UserID          gym.ID
	GymID           gym.ID
	Rating          Rating
	Comment         *string
	UserDisplayName *string
}

// NewReview creates a new review with validation
func NewReview(userID, gymID gym.ID, rating Rating) (*Review, error) {
	if !rating.IsValid() {
		return nil, gym.NewDomainError(gym.ErrInvalidRating, "invalid_rating", "rating must be between 1 and 5")
	}

	review := &Review{
		UserID: userID,
		GymID:  gymID,
		Rating: rating,
	}

	if err := review.Validate(); err != nil {
		return nil, err
	}

	return review, nil
}

// Validate validates review data
func (r *Review) Validate() error {
	if r.UserID == 0 {
		return gym.NewDomainError(gym.ErrInvalidInput, "invalid_user_id", "user ID is required")
	}

	if r.GymID == 0 {
		return gym.NewDomainError(gym.ErrInvalidInput, "invalid_gym_id", "gym ID is required")
	}

	if !r.Rating.IsValid() {
		return gym.NewDomainError(gym.ErrInvalidRating, "invalid_rating", "rating must be between 1 and 5")
	}

	if r.Comment != nil && len(*r.Comment) > 1000 {
		return gym.NewDomainError(gym.ErrInvalidInput, "invalid_comment", "comment too long")
	}

	return nil
}

// SetComment sets the review comment
func (r *Review) SetComment(comment string) {
	trimmed := strings.TrimSpace(comment)
	if trimmed == "" {
		r.Comment = nil
	} else {
		r.Comment = &trimmed
	}
}
