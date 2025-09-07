// internal/domain/review/value.go
// 役割: レビュードメインのバリューオブジェクト
// Rating(1..5保証)、PhotoURL等のレビュー関連バリューオブジェクトの定義と検証
package review

import (
	"database/sql/driver"
	"encoding/json"
	"net/url"
	"strings"
)

// Rating represents a review rating value object (1-5)
type Rating int

const (
	RatingMin Rating = 1
	RatingMax Rating = 5
)

// NewRating creates a new rating with validation
func NewRating(value int) (Rating, error) {
	rating := Rating(value)
	if !rating.IsValid() {
		return 0, &ValueError{Field: "rating", Message: "rating must be between 1 and 5"}
	}
	return rating, nil
}

// IsValid validates the rating value
func (r Rating) IsValid() bool {
	return r >= RatingMin && r <= RatingMax
}

// Int returns the int value of the rating
func (r Rating) Int() int {
	return int(r)
}

// PhotoURL represents a photo URL value object
type PhotoURL struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// NewPhotoURL creates a new photo URL with validation
func NewPhotoURL(urlStr, description string) (PhotoURL, error) {
	if urlStr == "" {
		return PhotoURL{}, &ValueError{Field: "url", Message: "URL cannot be empty"}
	}

	// Validate URL format
	if _, err := url.ParseRequestURI(urlStr); err != nil {
		return PhotoURL{}, &ValueError{Field: "url", Message: "invalid URL format"}
	}

	return PhotoURL{
		URL:         strings.TrimSpace(urlStr),
		Description: strings.TrimSpace(description),
	}, nil
}

// IsValid validates the photo URL
func (p PhotoURL) IsValid() bool {
	if p.URL == "" {
		return false
	}
	_, err := url.ParseRequestURI(p.URL)
	return err == nil
}

// PhotoURLs represents a slice of photo URLs stored as JSON
type PhotoURLs []PhotoURL

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
		return &ValueError{Field: "photos", Message: "cannot scan PhotoURLs"}
	}

	return json.Unmarshal(bytes, p)
}

// ValueError represents value object validation error
type ValueError struct {
	Field   string
	Message string
}

func (e *ValueError) Error() string {
	return e.Field + ": " + e.Message
}