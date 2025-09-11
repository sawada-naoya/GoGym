package model

import "time"

// Location represents gym location with geocoding
type Location struct {
	ID        int64     `json:"id"`
	Address   string    `json:"address"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	GymID     int64     `json:"gym_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (TODO: implement)
	Gym *Gym `json:"gym,omitempty"`
}