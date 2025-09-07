// internal/domain/review/entity.go
// 役割: レビュードメインのEntity/VO（Domain Layer）
// ビジネスルールと不変条件を持つ純粋なドメインオブジェクト。GORM/JSONタグは一切なし
package review

import (
	"gogym-api/internal/domain/common"
	"strings"
)

// Review represents a gym review aggregate root
type Review struct {
	common.BaseEntity
	UserID          common.ID
	GymID           common.ID
	Rating          Rating
	Comment         *string
	Photos          PhotoURLs
	UserDisplayName *string
}


// NewReview creates a new review with validation
func NewReview(userID, gymID common.ID, rating Rating) (*Review, error) {
	if !rating.IsValid() {
		return nil, common.NewDomainError(common.ErrInvalidRating, "invalid_rating", "rating must be between 1 and 5")
	}

	review := &Review{
		UserID: userID,
		GymID:  gymID,
		Rating: rating,
		Photos: make(PhotoURLs, 0),
	}

	if err := review.Validate(); err != nil {
		return nil, err
	}

	return review, nil
}

// Validate validates review data
func (r *Review) Validate() error {
	if r.UserID == 0 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_user_id", "user ID is required")
	}

	if r.GymID == 0 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_gym_id", "gym ID is required")
	}

	if !r.Rating.IsValid() {
		return common.NewDomainError(common.ErrInvalidRating, "invalid_rating", "rating must be between 1 and 5")
	}

	if r.Comment != nil && len(*r.Comment) > 1000 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_comment", "comment too long")
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

// AddPhoto adds a photo URL to the review
func (r *Review) AddPhoto(url PhotoURL) error {
	if !url.IsValid() {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_photo_url", "invalid photo URL")
	}

	// Initialize Photos if nil
	if r.Photos == nil {
		r.Photos = make(PhotoURLs, 0)
	}

	// Check for duplicates
	for _, existing := range r.Photos {
		if existing.URL == url.URL {
			return common.NewDomainError(common.ErrAlreadyExists, "photo_exists", "photo already exists in review")
		}
	}

	r.Photos = append(r.Photos, url)
	return nil
}

// RemovePhoto removes a photo URL from the review
func (r *Review) RemovePhoto(url string) {
	if r.Photos == nil {
		return
	}

	for i, existing := range r.Photos {
		if existing.URL == url {
			r.Photos = append(r.Photos[:i], r.Photos[i+1:]...)
			return
		}
	}
}

// HasPhotos returns true if the review has photos
func (r *Review) HasPhotos() bool {
	return r.Photos != nil && len(r.Photos) > 0
}