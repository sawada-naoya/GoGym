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
		return nil, NewDomainError("invalid_rating")
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
		return NewDomainError("invalid_user_id")
	}

	if r.GymID == 0 {
		return NewDomainError("invalid_gym_id")
	}

	if !r.Rating.IsValid() {
		return NewDomainError("invalid_rating")
	}

	if r.Comment != nil && len(*r.Comment) > 1000 {
		return NewDomainError("invalid_comment")
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
