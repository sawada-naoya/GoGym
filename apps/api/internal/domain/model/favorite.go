package model

import "time"

// Favorite represents a user's favorite gym
type Favorite struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	GymID     int64     `json:"gym_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (TODO: implement)
	User *User `json:"user,omitempty"`
	Gym  *Gym  `json:"gym,omitempty"`
}