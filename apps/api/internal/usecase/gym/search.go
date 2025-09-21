package gym

import (
	"context"
	"gogym-api/internal/domain/gym"
)

// SearchGyms searches gyms based on criteria
func (gu *GymUseCase) SearchGyms(ctx context.Context, req SearchGymRequest) (*SearchGymsResponse, error) {
	gu.logger.InfoContext(ctx, "searching gyms",
		"query", req.Query,
		"location", req.Location,
		"radius_m", req.RadiusM,
		"limit", req.Limit,
	)

	// Set default limit
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}

	// Set default radius
	if req.RadiusM != nil && (*req.RadiusM < 100 || *req.RadiusM > 50000) {
		defaultRadius := 5000
		req.RadiusM = &defaultRadius
	}

	searchQuery := gym.SearchQuery{
		Query:    req.Query,
		Location: req.Location,
		RadiusM:  req.RadiusM,
		Pagination: gym.Pagination{
			Cursor: req.Cursor,
			Limit:  req.Limit,
		},
	}

	result, err := gu.gymRepo.Search(ctx, searchQuery)
	if err != nil {
		gu.logger.ErrorContext(ctx, "failed to search gyms", "error", err)
		return nil, gym.NewDomainErrorWithCause(err, "search_failed", "failed to search gyms")
	}

	return &SearchGymsResponse{
		Gyms:       result.Items,
		NextCursor: &result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}