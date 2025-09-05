package common

import "time"

// ID represents a domain entity identifier
type ID int64

// Location represents geographic coordinates
type Location struct {
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
}

// IsValid validates location coordinates
func (l Location) IsValid() bool {
	return l.Latitude >= -90 && l.Latitude <= 90 &&
		l.Longitude >= -180 && l.Longitude <= 180
}

// Pagination represents pagination parameters
type Pagination struct {
	Cursor string
	Limit  int
}

// PaginatedResult represents paginated query result
type PaginatedResult[T any] struct {
	Items      []T
	NextCursor *string
	HasMore    bool
}

// BaseEntity represents common entity fields
type BaseEntity struct {
	ID        ID        `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// SearchQuery represents common search parameters
type SearchQuery struct {
	Query    string
	Location *Location
	RadiusM  *int
	Pagination
}