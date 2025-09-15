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

	gymIDs := make([]gym.ID, len(gyms))
	for i, g := range gyms {
		gymIDs[i] = g.ID
	}

	reviewStats, err := uc.gymRepo.GetReviewStatsForGyms(ctx, gymIDs)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to get review stats", "error", err)
	} else {
		for i := range gyms {
			if stats, exists := reviewStats[gyms[i].ID]; exists {
				gyms[i].AverageRating = stats.AverageRating
				gyms[i].ReviewCount = stats.ReviewCount
			}
		}
	}

	sort.Slice(gyms, func(i, j int) bool {
		if gyms[i].AverageRating != nil && gyms[j].AverageRating != nil {
			if *gyms[i].AverageRating != *gyms[j].AverageRating {
				return *gyms[i].AverageRating > *gyms[j].AverageRating
			}
			return gyms[i].ReviewCount > gyms[j].ReviewCount
		}
		if gyms[i].AverageRating != nil && gyms[j].AverageRating == nil {
			return true
		}
		if gyms[i].AverageRating == nil && gyms[j].AverageRating != nil {
			return false
		}
		return gyms[i].ID < gyms[j].ID
	})

	return &RecommendGymsResponse{
		Gyms:       gyms,
		NextCursor: &result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}
