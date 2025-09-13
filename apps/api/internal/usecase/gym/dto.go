package gym

import (
	"gogym-api/internal/domain/gym"
)

// SearchGymRequest represents search gym input
type SearchGymRequest struct {
	Query      string
	Location   *gym.Location
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
	Location    gym.Location
	Address     string
	City        *string
	Prefecture  *string
	PostalCode  *string
	TagNames    []string
}

// RecommendGymRequest represents recommend gym input
type RecommendGymRequest struct {
	UserLocation *gym.Location
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
	UserID gym.ID
	GymID  gym.ID
}

// GetFavoriteGymsRequest represents get favorite gyms input
type GetFavoriteGymsRequest struct {
	UserID gym.ID
	Limit  int
	Cursor string
}

// GetFavoriteGymsResponse represents get favorite gyms output
type GetFavoriteGymsResponse struct {
	Gyms       []gym.Gym
	NextCursor *string
	HasMore    bool
}

// GymDetailResponse represents detailed gym information for detail page
type GymDetailResponse struct {
	Gym gym.Gym
	// Reviews    []gym.Review   // TODO: Add when review functionality is implemented
	// Hours      []gym.Hour     // TODO: Add operating hours
	// Amenities  []gym.Amenity  // TODO: Add amenities/facilities
	// Images     []string       // Included in gym.Images field for now
}