package gym

import (
	"context"
	"gogym-api/internal/domain/common"
	"sort"
)

// RecommendGyms returns recommended gyms based on various criteria
func (uc *UseCase) RecommendGyms(ctx context.Context, req RecommendGymRequest) (*RecommendGymsResponse, error) {
	uc.logger.InfoContext(ctx, "getting recommended gyms",
		"user_location", req.UserLocation,
		"limit", req.Limit,
	)

	// Set default limit
	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 10
	}

	// Create search query for recommended gyms
	searchQuery := common.SearchQuery{
		Query:    "", // No text search for recommendations
		Location: req.UserLocation,
		RadiusM:  nil, // No radius limit for recommendations
		Pagination: common.Pagination{
			Cursor: req.Cursor,
			Limit:  req.Limit,
		},
	}

	result, err := uc.gymRepo.Search(ctx, searchQuery)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to get recommended gyms", "error", err)
		return nil, common.NewDomainErrorWithCause(err, "recommend_failed", "failed to get recommended gyms")
	}

	// Apply recommendation logic (sort by rating and reviews)
	gyms := result.Items
	
	// Add dummy ratings for gyms that don't have ratings (temporary until real data exists)
	for i := range gyms {
		if gyms[i].AverageRating == nil {
			// Generate consistent rating based on gym ID
			rating := 3.5 + float32((gyms[i].ID%10))*0.15
			gyms[i].AverageRating = &rating
		}
		if gyms[i].ReviewCount == 0 {
			// Generate consistent review count based on gym ID
			gyms[i].ReviewCount = int(gyms[i].ID%50) + 10
		}
	}

	// Sort by rating (highest first), then by review count
	sort.Slice(gyms, func(i, j int) bool {
		ratingI := float32(0)
		if gyms[i].AverageRating != nil {
			ratingI = *gyms[i].AverageRating
		}
		ratingJ := float32(0)
		if gyms[j].AverageRating != nil {
			ratingJ = *gyms[j].AverageRating
		}
		
		if ratingI == ratingJ {
			return gyms[i].ReviewCount > gyms[j].ReviewCount
		}
		return ratingI > ratingJ
	})

	return &RecommendGymsResponse{
		Gyms:       gyms,
		NextCursor: &result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}