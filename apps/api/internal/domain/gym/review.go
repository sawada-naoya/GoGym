package gym

import (
	"database/sql/driver"
	"encoding/json"
	"gogym-api/internal/domain/common"
	"strings"
)

// Review represents a gym review entity
type Review struct {
	common.BaseEntity
	UserID          common.ID       `json:"user_id" gorm:"not null;index;uniqueIndex:unique_user_gym_review,priority:1"`
	GymID           common.ID       `json:"gym_id" gorm:"not null;index;uniqueIndex:unique_user_gym_review,priority:2"`
	Rating          int             `json:"rating" gorm:"not null;index" validate:"required,min=1,max=5"`
	Comment         *string         `json:"comment" gorm:"type:text"`
	Photos          PhotoURLs       `json:"photos" gorm:"type:json"`
	UserDisplayName *string         `json:"user_display_name" gorm:"-"`
}

// TableName returns the table name for GORM
func (Review) TableName() string {
	return "reviews"
}

// PhotoURLs represents a slice of photo URLs stored as JSON
type PhotoURLs []string

// Value implements the driver.Valuer interface for database storage
func (p PhotoURLs) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

// Scan implements the sql.Scanner interface for database retrieval
func (p *PhotoURLs) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return common.NewDomainError(common.ErrInternal, "scan_error", "cannot scan PhotoURLs")
	}

	return json.Unmarshal(bytes, p)
}

// NewReview creates a new review with validation
func NewReview(userID, gymID common.ID, rating int) (*Review, error) {
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

	if r.Rating < 1 || r.Rating > 5 {
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
func (r *Review) AddPhoto(url string) error {
	if url == "" {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_photo_url", "photo URL cannot be empty")
	}

	// Initialize Photos if nil
	if r.Photos == nil {
		r.Photos = make(PhotoURLs, 0)
	}

	// Check for duplicates
	for _, existing := range r.Photos {
		if existing == url {
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
		if existing == url {
			r.Photos = append(r.Photos[:i], r.Photos[i+1:]...)
			return
		}
	}
}

// HasPhotos returns true if the review has photos
func (r *Review) HasPhotos() bool {
	return r.Photos != nil && len(r.Photos) > 0
}

// Favorite represents a user's favorite gym
type Favorite struct {
	common.BaseEntity
	UserID common.ID `json:"user_id" gorm:"not null;index;uniqueIndex:unique_user_gym_favorite,priority:1"`
	GymID  common.ID `json:"gym_id" gorm:"not null;index;uniqueIndex:unique_user_gym_favorite,priority:2"`
}

// TableName returns the table name for GORM
func (Favorite) TableName() string {
	return "favorites"
}

// NewFavorite creates a new favorite with validation
func NewFavorite(userID, gymID common.ID) (*Favorite, error) {
	favorite := &Favorite{
		UserID: userID,
		GymID:  gymID,
	}

	if err := favorite.Validate(); err != nil {
		return nil, err
	}

	return favorite, nil
}

// Validate validates favorite data
func (f *Favorite) Validate() error {
	if f.UserID == 0 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_user_id", "user ID is required")
	}

	if f.GymID == 0 {
		return common.NewDomainError(common.ErrInvalidInput, "invalid_gym_id", "gym ID is required")
	}

	return nil
}