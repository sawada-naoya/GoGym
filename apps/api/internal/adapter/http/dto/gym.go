package dto

import (
	"gogym-api/internal/domain/gym"
)

type GymResponse struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Location      string   `json:"location"`
	AverageRating *float32 `json:"average_rating,omitempty"`
	ReviewCount   int      `json:"review_count"`
	Tags          []string `json:"tags"`
}

// ToGymResponse converts domain Gym to DTO GymResponse
func ToGymResponse(g gym.Gym) GymResponse {
	description := ""
	if g.Description != nil {
		description = *g.Description
	}

	var tagNames []string
	for _, tag := range g.Tags {
		tagNames = append(tagNames, tag.Name)
	}

	return GymResponse{
		ID:            int(g.ID),
		Name:          g.Name,
		Description:   description,
		Location:      g.Location.String(),
		AverageRating: g.AverageRating,
		ReviewCount:   g.ReviewCount,
		Tags:          tagNames,
	}
}