package dto

import (
	"gogym-api/internal/domain/review"
	"time"
)

// ReviewResponse represents a gym review response
type ReviewResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Rating    int       `json:"rating"` // 1-5 rating
	GymID     int64     `json:"gym_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User *UserResponse `json:"user,omitempty"`
	Gym  *GymResponse  `json:"gym,omitempty"`
}

// ReviewListResponse represents a list of reviews
type ReviewListResponse struct {
	Reviews    []ReviewResponse `json:"reviews"`
	NextCursor *string          `json:"next_cursor"`
}

// FromReviewEntity converts review entity to DTO
func FromReviewEntity(r *review.Review) ReviewResponse {
	var content string
	if r.Comment != nil {
		content = *r.Comment
	}

	return ReviewResponse{
		ID:        int64(r.ID),
		Title:     "",
		Content:   content,
		Rating:    r.Rating.Int(),
		GymID:     int64(r.GymID),
		UserID:    int64(r.UserID),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
