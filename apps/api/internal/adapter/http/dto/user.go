package dto

import "time"

// UserResponse represents a system user response
type UserResponse struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations (TODO: implement)
	Reviews   []ReviewResponse   `json:"reviews,omitempty"`
	Favorites []FavoriteResponse `json:"favorites,omitempty"`
}