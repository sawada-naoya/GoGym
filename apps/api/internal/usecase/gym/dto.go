package gym

import (
	"gogym-api/internal/domain/common"
	"gogym-api/internal/domain/gym"
)

// SearchGymRequest represents search gym input
type SearchGymRequest struct {
	Query      string
	Location   *common.Location
	RadiusM    *int
	Cursor     string
	Limit      int
}

// SearchGymsResponse represents search gym output
type SearchGymsResponse struct {
	Gyms       []gym.Gym
	NextCursor *string
	HasMore    bool
}

// CreateGymRequest represents create gym input
type CreateGymRequest struct {
	Name        string
	Description *string
	Location    common.Location
	Address     string
	City        *string
	Prefecture  *string
	PostalCode  *string
	TagNames    []string
}

// RecommendGymRequest represents recommend gym input
type RecommendGymRequest struct {
	UserLocation *common.Location
	Limit        int
	Cursor       string
}

// RecommendGymsResponse represents recommend gym output
type RecommendGymsResponse struct {
	Gyms       []gym.Gym
	NextCursor *string
	HasMore    bool
}

// FavoriteGymRequest represents favorite gym input
type FavoriteGymRequest struct {
	UserID common.ID
	GymID  common.ID
}

// GetFavoriteGymsRequest represents get favorite gyms input
type GetFavoriteGymsRequest struct {
	UserID common.ID
	Limit  int
	Cursor string
}

// GetFavoriteGymsResponse represents get favorite gyms output
type GetFavoriteGymsResponse struct {
	Gyms       []gym.Gym
	NextCursor *string
	HasMore    bool
}