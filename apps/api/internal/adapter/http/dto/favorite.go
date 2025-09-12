package dto

import "time"

// FavoriteResponse represents a user's favorite gym response
type FavoriteResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	GymID     int64     `json:"gym_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (TODO: implement)
	User *UserResponse `json:"user,omitempty"`
	Gym  *GymResponse  `json:"gym,omitempty"`
}