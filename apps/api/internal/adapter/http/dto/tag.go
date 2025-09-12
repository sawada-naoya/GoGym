package dto

import "time"

// TagResponseFull represents a gym tag/category response with relations
type TagResponseFull struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (TODO: implement)
	Gyms []GymResponse `json:"gyms,omitempty"`
}

// GymTag represents many-to-many relationship between gyms and tags
type GymTag struct {
	ID        int64     `json:"id"`
	GymID     int64     `json:"gym_id"`
	TagID     int64     `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}