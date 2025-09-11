package model

import "time"

// User represents a system user
type User struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	CryptedPassword string   `json:"-"` // Hidden from JSON
	Salt           string    `json:"-"` // Hidden from JSON
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations (TODO: implement)
	Reviews   []Review   `json:"reviews,omitempty"`
	Favorites []Favorite `json:"favorites,omitempty"`
	Gyms      []Gym      `json:"gyms,omitempty"` // User-owned gyms
}