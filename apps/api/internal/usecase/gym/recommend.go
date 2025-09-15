package gym

import (
	"context"
	"gogym-api/internal/domain/gym"
	"sort"
)

// RecommendGyms returns recommended gyms based on various criteria
func (uc *UseCase) RecommendGyms(ctx context.Context, req RecommendGymRequest) (*RecommendGymsResponse, error) {
	uc.logger.InfoContext(ctx, "getting recommended gyms",
		"user_location", req.UserLocation,
		"limit", req.Limit,
	)

	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 10
	}

	searchQuery := gym.SearchQuery{
		Query:    "",
		Location: req.UserLocation,
		RadiusM:  nil,
		Pagination: gym.Pagination{
			Cursor: req.Cursor,
			Limit:  req.Limit,
		},
	}

	result, err := uc.gymRepo.Search(ctx, searchQuery)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to get recommended gyms", "error", err)
		return nil, err
	}

	gyms := result.Items

	// Sort gyms by ID (or could be modified to sort by other criteria later)
	sort.Slice(gyms, func(i, j int) bool {
		return gyms[i].ID < gyms[j].ID
	})

	return &RecommendGymsResponse{
		Gyms:       gyms,
		NextCursor: &result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}
