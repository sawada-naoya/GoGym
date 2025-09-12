package dto

import "time"

// UserResponse represents a system user response
type UserResponse struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	CryptedPassword string   `json:"-"` // Hidden from JSON
	Salt           string    `json:"-"` // Hidden from JSON
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations (TODO: implement)
	Reviews   []ReviewResponse   `json:"reviews,omitempty"`
	Favorites []FavoriteResponse `json:"favorites,omitempty"`
	Gyms      []GymResponse      `json:"gyms,omitempty"` // User-owned gyms
}