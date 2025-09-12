package dto

import "time"

// ReviewResponse represents a gym review response
type ReviewResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Rating    int       `json:"rating"`    // 1-5 rating
	ImageURL  *string   `json:"image_url"` // Optional image
	GymID     int64     `json:"gym_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations (TODO: implement)
	User *UserResponse `json:"user,omitempty"`
	Gym  *GymResponse  `json:"gym,omitempty"` // 循環参照回避のため、必要に応じてGymSummaryResponse等を検討
}