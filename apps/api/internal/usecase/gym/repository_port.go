package gym

import (
	"context"
	dom "gogym-api/internal/domain/gym"
)

// usecase → repository
type Repository interface {
	FindByID(ctx context.Context, id dom.ID) (*dom.Gym, error)
	Search(ctx context.Context, query dom.SearchQuery) (*dom.PaginatedResult[dom.Gym], error)
	GetReviewStatsForGyms(ctx context.Context, gymIDs []dom.ID) (map[dom.ID]*ReviewStats, error)
}

type ReviewStats struct {
	AverageRating *float32
	ReviewCount   int
}
