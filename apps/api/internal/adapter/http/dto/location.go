package dto

import "time"

// LocationResponse represents gym location with geocoding response
type LocationResponse struct {
	ID        int64     `json:"id"`
	Address   string    `json:"address"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	GymID     int64     `json:"gym_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (TODO: implement)
	Gym *GymResponse `json:"gym,omitempty"`
}