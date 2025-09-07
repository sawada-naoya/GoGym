// internal/domain/common/types.go
// 役割: 共通ドメイン型定義
// 全ドメインで共有される基本型とバリューオブジェクトの定義
package common

import (
	"fmt"
	"time"
)

// ID represents a unique identifier for entities
type ID int64

// BaseEntity provides common fields for all entities
type BaseEntity struct {
	ID        ID        `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Location represents geographical coordinates
type Location struct {
	Latitude  float64 `json:"latitude" gorm:"column:latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" gorm:"column:longitude" validate:"required,longitude"`
}

// IsValid validates location coordinates
func (l Location) IsValid() bool {
	return l.Latitude >= -90 && l.Latitude <= 90 &&
		l.Longitude >= -180 && l.Longitude <= 180
}

// String returns string representation of location
func (l Location) String() string {
	return fmt.Sprintf("(%f, %f)", l.Latitude, l.Longitude)
}

// SearchQuery represents search parameters
type SearchQuery struct {
	Query      string
	Location   *Location
	RadiusM    *int
	Pagination Pagination
}

// Pagination represents pagination parameters
type Pagination struct {
	Cursor string
	Limit  int
}

// PaginatedResult represents paginated response
type PaginatedResult[T any] struct {
	Items      []T
	NextCursor string
	HasMore    bool
}