package review

import (
	"context"
	dom "gogym-api/internal/domain/review"
)

// usecase → output (repository, external services, etc.)
type Repository interface {
	GetByGymID(ctx context.Context, gymID int64, cursor string, limit int) ([]dom.Review, string, error)
}
